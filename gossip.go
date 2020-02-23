package main

import (
	"context"
	"time"

	h "github.com/hashicorp/hyparview"
	"github.com/hashicorp/hyparview-example/proto"
)

type gossip struct {
	MaxHeat int   // config value
	Value   int32 // final value we got
	Hot     int   // gossip hotness
	Hops    int32 // number of hops the current value took to get here
	Seen    int   // if app == appSeen, we got every message
	Waste   int32 // count of app messages that didn't change the value
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

		peer := c.getPeer()
		if peer == nil {
			continue
		}

		hot, err := c.gossipSendSync(peer)
		if err != nil {
			c.failActive(peer)
		}

		if !hot {
			c.app.Hot -= 1
		}
	}
}

func (c *client) gossipSendSync(peer *h.Node) (bool, error) {
	cn, err := c.dial(peer)
	if err != nil {
		return false, err
	}
	grpc := cn.g

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := &proto.GossipRequest{Payload: c.app.Value, Hops: c.app.Hops}
	r, err := grpc.Gossip(ctx, req)
	return r.GetHot(), err
}

func (app *gossip) gossipRecv(payload, hops int32) bool {
	if app.Value >= payload {
		// Lifetime total count
		app.Waste += 1
		return false
	}
	app.Value = payload
	app.Hops = hops + 1
	app.Hot = app.MaxHeat

	// Lifetime total count
	app.Seen += 1
	return true
}
