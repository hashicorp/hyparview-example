package main

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"

	h "github.com/hashicorp/hyparview"
)

type ClientConfig struct {
	ID          string
	Addr        string
	Boot        string
	RootPEMFile string
}

type Gossip struct {
	Value int // final value we got
	Hot   int // gossip hotness
	Seen  int // if app == appSeen, we got every message
	Waste int // count of app messages that didn't change the value
}

type Client struct {
	config *ClientConfig
	tls    *tls.Config
	hv     *h.Hyparview
	app    *Gossip
	conn   map[string]*tls.Conn
	in     []h.Message
	out    []h.Message
}

func makeID() string {
	bs := make([]byte, 8)
	rand.Read(bs)
	return fmt.Sprintf("%160x", bs)
}

func MakeClient(c *ClientConfig) (*Client, error) {
	addr := fmt.Sprintf("%s:%d", c.Addr, c.Port)
	cert, err := ioutil.ReadFile(c.RootPEMFile)
	if err != nil {
		return nil, err
	}

	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM(cert)
	if !ok {
		return nil, fmt.Errorf("failed to parse root cert %s", c.RootPEMFile)
	}

	return &Client{
		config: c,
		tls:    &tls.Config{RootCAs: roots},
		hy:     h.CreateView(&h.Node{ID: id, Addr: addr}),
		app:    &Gossip{},
	}
}

func (c *Client) Dial(node *h.Node) (tls.Conn, error) {
	cn, ok := c.conn[node.Addr]
	if ok {
		return cn, nil
	}

	cn, err := tls.Dial("tcp", node.Addr, c.tls)
	if err == nil {
		return nil, err
	}

	c.conn[node.Addr] = cn
	return cn, nil
}

func (c *Client) Outbox(ms ...h.Message) {
	c.out = append(c.out, ms...)
}

func (c *Client) Send(m *h.Message) (resp *h.NeighborRefuse, err error) {
	bs, err := c.sendSync(m.To(), json.Marshal(m))
	if len(bs) > 0 {
		err = json.Unmarshal(bs, resp)
	}
	return resp, err
}

func (c *Client) sendSync(peer *h.Node, bs []byte) ([]byte, error) {
	conn, err := c.Dial(peer.Addr)
	if err != nil {
		return nil, err
	}

	len, err := conn.Write(bs)

	if err != nil {
		c.failActive(peer)
		return nil, err
	}

	if len(bs) != len {
		c.failActive(peer)
		return nil, fmt.Errorf("only wrote %d of %d bytes", len, len(bs))
	}

	var resp []byte
	err = conn.Read(resp)
	return resp, err
}

func (c *Client) Drain(count int) {
	for i := count; i > 0; i-- {
		ms, err := c.Send(c.out[0])
		if err != nil {
			continue
		}
		if len(ms) > 0 {
			c.Outbox(ms...)
		}
		c.out = c.out[1:]
	}
}

func (c *Client) failActive(peer *h.Node) error {
	v := c.hv
	idx := v.Active.ContainsIndex(peer)
	if idx < 0 {
		return nil
	}

	for _, n := range v.Passive.Shuffled() {
		if v.Active.IsEmpty() {
			// High priority can't be rejected, so send async
			m := h.SendNeighbor(n, v.Self, h.HighPriority)
			c.Outbox(m)
			break
		} else {
			m := h.SendNeighbor(n, v.Self, h.LowPriority)
			resp, err := c.Send(m)
			// Either moved to the active view, or failed
			v.DelPassive(n)
			// empty low priority response is success
			if resp == nil && err == nil {
				c.Outbox(v.AddActive(n)...)
				break
			}
		}
	}
	return nil
}
