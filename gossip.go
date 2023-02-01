// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"time"

	"github.com/hashicorp/hyparview-example/proto"
)

type gossip struct {
	Payload int32 // final value we got
	Hops    int32 // number of hops the current value took to get here
	Seen    int   // if app == appSeen, we got every message
	Waste   int32 // count of app messages that didn't change the value
}

func newGossip(maxHeat int) *gossip {
	return &gossip{}
}

func (c *client) gossipStart() {
	c.app.Payload = c.app.Payload + 1
	c.app.Hops = 0
	c.gossipForward(&proto.GossipRequest{
		Payload: c.app.Payload,
		Hops:    1,
		From:    c.hv.Active.RandNode().Addr,
	})
}

func (c *client) gossipForward(m *proto.GossipRequest) {
	from := node(m.From)
	for _, peer := range c.hv.Active.Shuffled() {
		if peer.Equal(from) {
			continue
		}

		cn, err := c.dial(peer)
		if err != nil {
			// failActive is async
			c.failActive(peer)
			continue
		}
		grpc := cn.g

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		_, err = grpc.Gossip(ctx, m)
	}
}

func (c *client) gossipRecv(m *proto.GossipRequest) {
	if c.app.Payload >= m.Payload {
		// Lifetime total count
		c.app.Waste += 1
		return
	}
	c.app.Payload = m.Payload
	m.Hops += 1
	c.app.Hops = m.Hops

	// Lifetime total count
	c.app.Seen += 1

	c.gossipForward(m)
}
