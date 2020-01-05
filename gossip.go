package main

import (
	"context"
	"time"

	h "github.com/hashicorp/hyparview"
	"github.com/hashicorp/hyparview-example/proto"
)

type gossip struct {
	MaxHeat int // config value
	Value   int // final value we got
	Hot     int // gossip hotness
	Seen    int // if app == appSeen, we got every message
	Waste   int // count of app messages that didn't change the value
}

func newGossip(maxHeat int) *gossip {
	return &gossip{
		MaxHeat: maxHeat,
	}
}

func (c *client) gossipSend(payload int) {
	for c.app.Hot > 0 {
		if h.Rint(c.app.MaxHeat) > c.app.Hot {
			continue
		}

		peer := c.hv.Peer()
		if peer == nil {
			c.failActive(nil) // ignore errors
			continue
		}

		hot, err := c.gossipSendSync(peer, payload)
		if err != nil {
			c.failActive(peer)
		}

		if !hot {
			c.app.Hot -= 1
		}
	}
}

func (c *client) gossipSendSync(peer *h.Node, payload int) (bool, error) {
	cn, err := c.dial(peer)
	if err != nil {
		return false, err
	}
	grpc := cn.g

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := &proto.GossipRequest{Payload: int32(payload)}
	r, err := grpc.Gossip(ctx, req)
	return r.GetHot(), err
}

func (app *gossip) gossipRecv(payload int) bool {
	if app.Value >= payload {
		app.Waste += 1
		return false
	}
	app.Value = payload
	app.Seen += 1
	app.Hot = app.MaxHeat
	return true
}
