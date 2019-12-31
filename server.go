//go:generate protoc -I . --go_out=plugins=grpc:./proto ./proto/hyparview.proto ./proto/gossip.proto

package main

import (
	"context"

	h "github.com/hashicorp/hyparview"
	"github.com/hashicorp/hyparview-example/proto/proto"
)

type server struct {
	c *client
}

func newServer(c *client) *server {
	return &server{c: c}
}

func (s *server) Gossip(ctx context.Context, req *proto.GossipRequest) (*proto.GossipResponse, error) {
	hot := s.c.app.gossipRecv(int(req.Payload))
	return &proto.GossipResponse{Hot: hot}, nil
}

func (s *server) Join(ctx context.Context, req *proto.FromRequest) (*proto.Empty, error) {
	to, from := s.c.hv.Self, &h.Node{Addr: req.From}
	ms := s.c.hv.RecvJoin(h.SendJoin(to, from))
	s.c.outbox(ms...)
	return &proto.Empty{}, nil
}

func (s *server) ForwardJoin(ctx context.Context, req *proto.ForwardJoinRequest) (*proto.Empty, error) {
	to, from := s.c.hv.Self, &h.Node{Addr: req.From}
	join := &h.Node{Addr: req.Join}
	ttl := int(req.Ttl)
	ms := s.c.hv.RecvForwardJoin(h.SendForwardJoin(to, from, join, ttl))
	s.c.outbox(ms...)
	return &proto.Empty{}, nil
}

func (s *server) Disconnect(ctx context.Context, req *proto.FromRequest) (*proto.Empty, error) {
	to, from := s.c.hv.Self, &h.Node{Addr: req.From}
	s.c.hv.RecvDisconnect(h.SendDisconnect(to, from))
	return &proto.Empty{}, nil
}

func (s *server) Neighbor(ctx context.Context, req *proto.NeighborRequest) (*proto.NeighborResponse, error) {
	to, from := s.c.hv.Self, &h.Node{Addr: req.From}
	priority := req.Priority
	ms := s.c.hv.RecvNeighbor(h.SendNeighbor(to, from, priority))
	accept := len(ms) == 0
	return &proto.NeighborResponse{Accept: accept}, nil
}

func (s *server) Shuffle(ctx context.Context, req *proto.ShuffleRequest) (*proto.ShuffleReply, error) {
	to, from := s.c.hv.Self, &h.Node{Addr: req.From}
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
