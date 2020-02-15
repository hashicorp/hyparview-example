package main

import (
	"context"
	"log"

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

func (s *server) View(ctx context.Context, req *proto.Empty) (*proto.ViewResponse, error) {
	return &proto.ViewResponse{
		From: s.c.hv.Self.Addr
		Active:  sliceNodeAddr(s.c.hv.Active.Nodes),
		Passive: sliceNodeAddr(s.c.hv.Passive.Nodes),
		App: s.c.app.Value,
		Hops: s.c.app.Hops,
	}
}

func (s *server) Join(ctx context.Context, req *proto.FromRequest) (*proto.Empty, error) {
	// log.Printf("info join recv: %s\n", req.From)
	to, from := s.c.hv.Self, node(req.From)
	ms := s.c.hv.RecvJoin(h.SendJoin(to, from))
	s.c.outbox(ms...)
	return &proto.Empty{}, nil
}

func (s *server) ForwardJoin(ctx context.Context, req *proto.ForwardJoinRequest) (*proto.Empty, error) {
	to, from := s.c.hv.Self, node(req.From)
	join := &h.Node{Addr: req.Join}
	ttl := int(req.Ttl)
	ms := s.c.hv.RecvForwardJoin(h.SendForwardJoin(to, from, join, ttl))
	s.c.outbox(ms...)
	return &proto.Empty{}, nil
}

func (s *server) Disconnect(ctx context.Context, req *proto.FromRequest) (*proto.Empty, error) {
	to, from := s.c.hv.Self, node(req.From)
	s.c.hv.RecvDisconnect(h.SendDisconnect(to, from))
	return &proto.Empty{}, nil
}

func (s *server) Neighbor(ctx context.Context, req *proto.NeighborRequest) (*proto.NeighborResponse, error) {
	to, from := s.c.hv.Self, node(req.From)
	priority := req.Priority
	ms := s.c.hv.RecvNeighbor(h.SendNeighbor(to, from, priority))
	accept := len(ms) == 0
	return &proto.NeighborResponse{Accept: accept}, nil
}

func (s *server) Shuffle(ctx context.Context, req *proto.ShuffleRequest) (*proto.ShuffleReply, error) {
	// log.Printf("info shuffle recv: %s\n", req.From)
	to, from := s.c.hv.Self, node(req.From)
	active := sliceAddrNode(req.Active)
	passive := sliceAddrNode(req.Passive)
	ttl := int(req.Ttl)
	ms := s.c.hv.RecvShuffle(h.SendShuffle(to, from, active, passive, ttl))

	var res *proto.ShuffleReply

	for _, m := range ms {
		switch v := m.(type) {
		case *h.ShuffleRequest:
			s.c.outbox(v)
		case *h.ShuffleReply:
			res = &proto.ShuffleReply{Passive: sliceNodeAddr(v.Passive)}
		}
	}

	return res, nil
}
