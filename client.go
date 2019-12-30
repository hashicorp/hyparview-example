package main

import (
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
	Port        int
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

func MakeClient(c *ClientConfig) (*Client, err) {
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

func (c *Client) GossipSend(payload int) {
	for c.app.Hot > 0 {
		if h.Rint(c.world.config.gossipHeat) > c.app.Hot {
			continue
		}

		peer := c.hv.Peer()
		if peer == nil {
			c.failActive(nil) // ignore errors
			continue
		}

		hot, err := c.gossipSend(peer, payload)
		if err != nil {
			c.failActive(peer)
		}

		if !hot {
			c.app.Hot -= 1
		}
	}
}

func (c *Client) gossipSend(peer *h.Node, payload int) (bool, error) {
	bs, err := c.sendSync(peer, []byte{payload})
	var out bool
	if len(bs) > 0 {
		out = bs[0] == 1
	}
	return out, err
}

func (c *Client) GossipRecv(payload int) bool {
	if c.app.Value >= payload {
		c.app.Waste += 1
		return false
	}
	c.app.Value = payload
	c.app.Seen += 1
	c.app.Hot = c.world.config.gossipHeat
	return true
}
