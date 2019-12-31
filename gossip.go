package main

func (c *Client) gossipSend(payload int) {
	for c.app.Hot > 0 {
		if h.Rint(c.world.config.gossipHeat) > c.app.Hot {
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

func (c *Client) gossipSendSync(peer *h.Node, payload int) (bool, error) {
	bs, err := c.sendSync(peer, []byte{payload})
	var out bool
	if len(bs) > 0 {
		out = bs[0] == 1
	}
	return out, err
}

func (c *Client) gossipRecv(payload int) bool {
	if c.app.Value >= payload {
		c.app.Waste += 1
		return false
	}
	c.app.Value = payload
	c.app.Seen += 1
	c.app.Hot = c.world.config.gossipHeat
	return true
}
