package main

import (
	"context"
	// "log"

	h "github.com/hashicorp/hyparview"
	"github.com/hashicorp/hyparview-example/proto"
)

type server struct {
	c *client
}

func newServer(c *client) *server {
	return &server{c: c}
}

func (s *server) Gossip(ctx context.Context, req *proto.GossipRequest) (*proto.GossipResponse, error) {
	hot := s.c.app.gossipRecv(req.Payload, req.Hops)
	// if hot {
	// 	log.Printf("info gossip recv: %d\n", req.Payload)
	// }

	return &proto.GossipResponse{Hot: hot}, nil
}

func (s *server) View(ctx context.Context, req *proto.StatEmpty) (*proto.ViewResponse, error) {
	return &proto.ViewResponse{
		From:    s.c.hv.Self.Addr,
		Active:  sliceNodeAddr(s.c.hv.Active.Nodes),
		Passive: sliceNodeAddr(s.c.hv.Passive.Nodes),
		App:     s.c.app.Value,
		Hops:    s.c.app.Hops,
		Waste:   s.c.app.Waste,
	}, nil
}

func (s *server) Join(ctx context.Context, req *proto.FromRequest) (*proto.HyparviewEmpty, error) {
	// log.Printf("info join recv: %s\n", req.From)
	to, from := s.c.hv.Self, node(req.From)
	s.c.inbox(h.SendJoin(to, from))
	return &proto.HyparviewEmpty{}, nil
}

func (s *server) ForwardJoin(ctx context.Context, req *proto.ForwardJoinRequest) (*proto.HyparviewEmpty, error) {
	to, from := s.c.hv.Self, node(req.From)
	join := &h.Node{Addr: req.Join}
	ttl := int(req.Ttl)
	s.c.inbox(h.SendForwardJoin(to, from, join, ttl))
	return &proto.HyparviewEmpty{}, nil
}

func (s *server) Disconnect(ctx context.Context, req *proto.FromRequest) (*proto.HyparviewEmpty, error) {
	to, from := s.c.hv.Self, node(req.From)
	s.c.inbox(h.SendDisconnect(to, from))
	return &proto.HyparviewEmpty{}, nil
}

func (s *server) Neighbor(ctx context.Context, req *proto.NeighborRequest) (*proto.NeighborResponse, error) {
	to, from := s.c.hv.Self, node(req.From)
	priority := req.Priority
	k := s.c.inboxAwait(h.SendNeighbor(to, from, priority))
	ms := <-k
	accept := ms == nil
	return &proto.NeighborResponse{Accept: accept}, nil
}

func (s *server) Shuffle(ctx context.Context, req *proto.ShuffleRequest) (*proto.HyparviewEmpty, error) {
	to, from := s.c.hv.Self, node(req.From)
	active := sliceAddrNode(req.Active)
	passive := sliceAddrNode(req.Passive)
	ttl := int(req.Ttl)
	s.c.inbox(h.SendShuffle(to, from, active, passive, ttl))
	return &proto.HyparviewEmpty{}, nil
}

func (s *server) ShuffleReply(ctx context.Context, req *proto.ShuffleReplyRequest) (*proto.HyparviewEmpty, error) {
	to, from := s.c.hv.Self, node(req.From)
	passive := sliceAddrNode(req.Passive)
	s.c.inbox(h.SendShuffleReply(to, from, passive))
	return &proto.HyparviewEmpty{}, nil
}
