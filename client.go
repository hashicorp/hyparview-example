package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"sync"
	"time"

	h "github.com/hashicorp/hyparview"
	"github.com/hashicorp/hyparview-example/proto"
	"google.golang.org/grpc"
)

type clientConfig struct {
	// sendFanOut count of sending threads
	sendFanOut int
	// shuffleSeconds max seconds between shuffle requests
	shuffleSeconds int
	hvClientCount  int
	hvInboxBuffer  int
	hvOutboxBuffer int
	gossipMaxHeat  int

	// statParseFanOut count stats deserialization threads
	statParseFanOut int
	// statParseBuffer buffer stats messages per parseFanOut process
	statParseBuffer  int
	statMillis       int
	statUpdateBuffer int
	statUDPBuffer    int

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

	app *gossip

	conn     map[string]*conn
	connLock sync.RWMutex

	hv *h.Hyparview
	// in serializes changes to hv
	in chan *message
	// out fans out to a dedicated set of senders
	out chan h.Message
}

func newID() string {
	bs := make([]byte, 8)
	rand.Read(bs)
	return fmt.Sprintf("%x", bs)
}

func newClient(c *clientConfig) *client {
	return &client{
		config: c,
		hv:     h.CreateView(&h.Node{ID: c.addr, Addr: c.addr}, c.hvClientCount),
		app:    newGossip(c.gossipMaxHeat),
		conn:   map[string]*conn{},
		in:     make(chan *message, c.hvInboxBuffer),
		out:    make(chan h.Message, c.hvOutboxBuffer),
	}
}

// dial connects to a client and caches the connection
func (c *client) dial(node *h.Node) (*conn, error) {
	c.connLock.RLock()
	cn, ok := c.conn[node.Addr]
	c.connLock.RUnlock()
	if ok {
		return cn, nil
	}

	cn, err := c.justDial(node)
	if err != nil {
		return nil, err
	}

	c.connLock.Lock()
	defer c.connLock.Unlock()
	c.conn[node.Addr] = cn
	return cn, nil
}

// justDial creates a temporary connection. split out for ShuffleReply
func (c *client) justDial(node *h.Node) (*conn, error) {
	// Client name must match the dns name of the server
	creds, err := clientCreds(c.config, "localhost")
	if err != nil {
		return nil, fmt.Errorf("error client credential: %v", err)
	}

	g, err := grpc.Dial(node.Addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, err
	}

	cn := &conn{
		c: g,
		h: proto.NewHyparviewClient(g),
		g: proto.NewGossipClient(g),
	}

	return cn, nil
}

func (c *client) drop(node *h.Node) {
	c.connLock.Lock()
	defer c.connLock.Unlock()
	cn, _ := c.conn[node.Addr]
	if cn != nil {
		cn.c.Close()
	}
	delete(c.conn, node.Addr)
}

func (c *client) send(m h.Message) (err error) {
	conn, err := c.dial(m.To())
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	switch v := m.(type) {
	case *h.JoinRequest:
		r := &proto.FromRequest{From: c.hv.Self.Addr}
		_, err = conn.h.Join(ctx, r)

	case *h.ForwardJoinRequest:
		r := &proto.ForwardJoinRequest{
			Ttl:  int32(v.TTL),
			Join: v.Join.Addr,
			From: v.From.Addr,
		}
		_, err = conn.h.ForwardJoin(ctx, r)

	case *h.DisconnectRequest:
		r := &proto.FromRequest{From: c.hv.Self.Addr}
		_, err = conn.h.Disconnect(ctx, r)
		c.drop(v.To())

	case *h.ShuffleRequest:
		r := &proto.ShuffleRequest{
			Ttl:     int32(v.TTL),
			From:    c.hv.Self.Addr,
			Active:  sliceNodeAddr(v.Active),
			Passive: sliceNodeAddr(v.Passive),
		}
		_, err = conn.h.Shuffle(ctx, r)

	case *h.ShuffleReply:
		// Only ShuffleReply is sent to a peer not in our active view via a
		// temporary conn. If we don't have the conn cached, we'll close it at the
		// end of this send
		if conn == nil {
			conn, err = c.justDial(m.To())
			if err != nil {
				return err
			}
			defer conn.c.Close()
		}
		r := &proto.ShuffleReplyRequest{
			From:    c.hv.Self.Addr,
			Passive: sliceNodeAddr(v.Passive),
		}
		_, err = conn.h.ShuffleReply(ctx, r)
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

func (c *client) outbox(ms ...h.Message) {
	for _, m := range ms {
		c.out <- m
	}
}

// message wraps the hyparview message with an optional return channel. For calls (like
// shuffle) that need a return value, recv will write the messages to the `k` channel.
// Regular calls will simply produce outbox messages as a side effect
type message struct {
	m h.Message
	k chan h.Message
	// fail is message type all by itself, we need to process failures in the recv
	// thread
	fail *h.Node
}

func (c *client) inbox(ms ...h.Message) {
	for _, m := range ms {
		c.in <- &message{m: m}
	}
}

// inboxAwait returns a channel which blocks until the response is available
func (c *client) inboxAwait(m h.Message) chan h.Message {
	k := make(chan h.Message)
	c.in <- &message{m: m, k: k}
	return k
}

// recv is the single threaded consumer of hyparview messages
func (c *client) recv(m *message) {
	// Dirty hack to process failActive in this thread
	if m.fail != nil {
		c.recvFailActive(m.fail)
		return
	}

	v := c.hv.Copy()

	switch m1 := m.m.(type) {
	case *h.NeighborRequest:
		for _, msg := range v.RecvNeighbor(m1) {
			m.k <- msg
		}
		close(m.k)
	case *h.JoinRequest:
		c.outbox(v.RecvJoin(m1)...)
	case *h.ForwardJoinRequest:
		c.outbox(v.RecvForwardJoin(m1)...)
	case *h.DisconnectRequest:
		v.RecvDisconnect(m1)
		c.drop(m1.From)
	case *h.ShuffleRequest:
		c.outbox(v.RecvShuffle(m1)...)
	case *h.ShuffleReply:
		v.RecvShuffleReply(m1)
	default:
		// log unimplemented?
	}

	c.hv = v
}

// getPeer returns a peer if possible, and maintains the active view
func (c *client) getPeer() *h.Node {
	// Add to the active view so we're always mostly full, but one message at a time
	if !c.hv.Active.IsFull() {
		c.failActive(nil)
	}
	return c.hv.Peer()
}

func (c *client) failActive(peer *h.Node) {
	c.in <- &message{
		fail: peer,
	}
}

func (c *client) recvFailActive(peer *h.Node) {
	v := c.hv.Copy()

	// peer may be nil, which means that the client found our ActiveView empty. In that
	// case we can't drop anything, but should recover our active view
	if peer != nil {
		idx := v.Active.ContainsIndex(peer)
		if idx >= 0 {
			v.Active.DelIndex(idx)
			c.drop(peer)
		}
	}

	for _, n := range v.Passive.Shuffled() {
		pri := h.LowPriority
		if v.Active.IsEmpty() {
			pri = h.HighPriority
		}

		// Send sync so we can detect errors
		m := h.SendNeighbor(n, v.Self, pri)
		res, err := c.sendNeighbor(m)

		// If refused, keep going, and keep this server in the list
		if pri == h.LowPriority && res != nil {
			// log.Printf("info: failActive refuse %s", n.Addr)
			continue
		}

		// Either moved to the active view, or failed
		v.DelPassive(n)

		if err != nil {
			// log.Printf("info: failActive error %s %v", n.Addr, err)
			continue
		}

		c.outbox(v.AddActive(n)...)
		break
	}

	c.hv = v
}

func (c *client) statSleep() {
	time.Sleep(time.Duration(h.Rint(c.config.statMillis)) * time.Millisecond)
}

func (c *client) shuffleSleep() {
	time.Sleep(time.Duration(h.Rint(c.config.shuffleSeconds)) * time.Second)
}
