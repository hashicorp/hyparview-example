package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"time"

	h "github.com/hashicorp/hyparview"
	"github.com/hashicorp/hyparview-example/proto"
	"google.golang.org/grpc"
)

type clientConfig struct {
	id         string
	addr       string
	bootstrap  string
	caCert     string
	serverCert string
	serverKey  string
	clientCert string
	clientKey  string
}

type conn struct {
	c *grpc.ClientConn
	h proto.HyparviewClient
	g proto.GossipClient
}

type client struct {
	config *clientConfig
	grpc   []grpc.DialOption
	hv     *h.Hyparview
	app    *gossip
	conn   map[string]*conn
	out    chan h.Message
}

func newID() string {
	bs := make([]byte, 8)
	rand.Read(bs)
	return fmt.Sprintf("%x", bs)
}

func newClient(c *clientConfig) *client {
	return &client{
		config: c,
		hv:     h.CreateView(&h.Node{ID: c.addr, Addr: c.addr}, 10000),
		app:    newGossip(4),
		conn:   map[string]*conn{},
		out:    make(chan h.Message, 2048),
	}
}

func (c *client) dial(node *h.Node) (*conn, error) {
	cn, ok := c.conn[node.Addr]
	if ok {
		return cn, nil
	}

	// Client name must match the dns name of the server
	creds, err := clientCreds(c.config, "localhost")
	if err != nil {
		return nil, fmt.Errorf("error client credential: %v", err)
	}

	g, err := grpc.Dial(node.Addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, err
	}

	cn = &conn{
		c: g,
		h: proto.NewHyparviewClient(g),
		g: proto.NewGossipClient(g),
	}

	c.conn[node.Addr] = cn
	return cn, nil
}

func (c *client) drop(node *h.Node) {
	cn, _ := c.conn[node.Addr]
	if cn != nil {
		cn.c.Close()
	}
	delete(c.conn, node.Addr)
}

func (c *client) send(m h.Message) (err error) {
	cn, err := c.dial(m.To())
	if err != nil {
		return err
	}
	grpc := cn.h

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	switch v := m.(type) {
	case *h.JoinRequest:
		r := &proto.FromRequest{From: c.hv.Self.Addr}
		_, err = grpc.Join(ctx, r)

	case *h.ForwardJoinRequest:
		r := &proto.ForwardJoinRequest{
			Ttl:  int32(v.TTL),
			Join: v.Join.Addr,
			From: v.From.Addr,
		}
		_, err = grpc.ForwardJoin(ctx, r)

	case *h.DisconnectRequest:
		r := &proto.FromRequest{From: c.hv.Self.Addr}
		_, err = grpc.Disconnect(ctx, r)

	case *h.NeighborRequest:
		// Only in `out` if high priority, safe to ignore the response
		r := &proto.NeighborRequest{Priority: v.Priority, From: v.From.Addr}
		_, err = grpc.Neighbor(ctx, r)
	}
	return err
}

func (c *client) sendNeighbor(m *h.NeighborRequest) (*h.NeighborRefuse, error) {
	cn, err := c.dial(m.To())
	if err != nil {
		return nil, err
	}
	grpc := cn.h

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := &proto.NeighborRequest{Priority: m.Priority, From: m.From.Addr}
	r, err := grpc.Neighbor(ctx, req)
	if err != nil {
		return nil, err
	}

	if r.Accept {
		return nil, nil
	}

	return h.SendNeighborRefuse(c.hv.Self, &h.Node{Addr: r.From}), nil
}

func (c *client) sendShuffle(m *h.ShuffleRequest) (res *h.ShuffleReply, err error) {
	to := m.To()
	if to == nil {
		return nil, nil
	}
	cn, err := c.dial(m.To())
	if err != nil {
		return nil, fmt.Errorf("dial: %v", err)
	}
	grpc := cn.h

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	req := &proto.ShuffleRequest{
		Ttl:     int32(m.TTL),
		Active:  sliceNodeAddr(m.Active),
		Passive: sliceNodeAddr(m.Passive),
		From:    m.From.Addr,
	}

	r, err := grpc.Shuffle(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("grpc: %v", err)
	}

	return h.SendShuffleReply(c.hv.Self, &h.Node{Addr: r.From}, sliceAddrNode(r.Passive)), nil
}

func (c *client) outbox(ms ...h.Message) {
	for _, m := range ms {
		c.out <- m
	}
}

func (c *client) failActive(peer *h.Node) {
	// peer may be nil, which means that the client found our ActiveView empty. In that
	// case we can't drop anything, but should recover our active view
	v := c.hv

	if peer != nil {
		idx := v.Active.ContainsIndex(peer)
		if idx > -1 {
			v.Active.DelIndex(idx)
			c.drop(peer)
		}
	}

	passive := v.Passive.Copy()
	for _, n := range passive.Shuffled() {
		pri := h.LowPriority
		if v.Active.IsEmpty() {
			pri = h.HighPriority
		}

		// Send sync so we can detect errors
		m := h.SendNeighbor(n, v.Self, pri)
		res, err := c.sendNeighbor(m)

		// If refused, keep going, and keep this server in the list
		if res != nil {
			log.Printf("info: failActive refuse %s", n.Addr)
			continue
		}

		// Either moved to the active view, or failed
		v.DelPassive(n)

		if err != nil {
			log.Printf("info: failActive error %s %v", n.Addr, err)
			continue
		}

		c.outbox(v.AddActive(n)...)
		break
	}
}
