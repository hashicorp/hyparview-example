package main

import (
	"log"
	"net"
)

type stats struct {
	safe map[string]*peerStat
}

func newStats() *stats {
	return &stats{
		safe: map[string]*peerStat{},
	}
}

// runStatServer runs forever, populating the stats struct
func runStatServer(addr string, stats *stats, c *client) {
	log.Printf("debug: starting stat server %s", addr)

	parsers := c.config.statParseFanOut
	inbox := make(chan []byte, parsers*c.config.statParseBuffer)
	update := make(chan *peerStat, c.config.statUpdateBuffer)

	// Start some parsers
	w := &wireStat{}
	for i := 0; i < parsers; i++ {
		go func() {
			for {
				p := <-inbox
				w.Parse(p[:])
				update <- w.peerStat()
			}
		}()
	}

	// Start the updater
	go func() {
		for {
			p := <-update
			stats.safe[p.From] = p
		}
	}()

	// Update our own stats
	go func() {
		for {
			c.statSleep()
			update <- c.wireStat().peerStat()
		}
	}()

	// Finally, one listener
	pc, _ := net.ListenPacket("udp", addr)
	defer pc.Close()

	for {
		buf := make([]byte, STAT_SIZE)
		pc.ReadFrom(buf)
		inbox <- buf
	}
}

// runStatClient blocks forever sending stats on a random interval to the remote addr
func runStatClient(c *client, addr string) {
	log.Printf("debug: stat client %s", addr)

	// Under the impression that udp dial won't return an error :shrug:
	conn, _ := net.Dial("udp", addr)
	defer conn.Close()

	for {
		c.statSleep()
		ws := c.wireStat()
		conn.Write(ws.Bytes())
	}
}
