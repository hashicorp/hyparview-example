package main

import (
	"net"
	"time"

	h "github.com/hashicorp/hyparview"
)

type stats struct {
	safe map[string]*peerStat
}

func newStats() *stats {
	return &stats{}
}

// runStatServer runs forever, populating the stats struct
func runStatServer(addr string, stats *stats) {
	parsers := 10
	inbox := make(chan []byte, parsers*2)
	update := make(chan *peerStat, 5)

	// Start some parsers
	w := &wireStat{}
	for i := 0; i < parsers; i++ {
		go func() {
			p := <-inbox
			w.Parse(p[:])
			update <- w.peerStat()
		}()
	}

	// Start the updater
	go func() {
		p := <-update
		stats.safe[p.From] = p
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
	conn, _ := net.Dial("udp", addr)
	defer conn.Close()

	for {
		time.Sleep(time.Duration(h.Rint(500)) * time.Millisecond)
		conn.Write(c.wireStat().Bytes())
	}
}
