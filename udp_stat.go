package main

import (
	"log"
	"net"
	"time"

	h "github.com/hashicorp/hyparview"
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

	parsers := 10
	inbox := make(chan []byte, parsers*2)
	update := make(chan *peerStat, 5)

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
			time.Sleep(time.Duration(h.Rint(500)) * time.Millisecond)
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

	// Retry connections until we get through
	conn, err := net.Dial("udp", addr)
	for {
		if err == nil {
			break
		}
		log.Printf("info: stat client %v", err)
		time.Sleep(time.Second)
		conn, err = net.Dial("udp", addr)
	}
	defer conn.Close()

	for {
		time.Sleep(time.Duration(h.Rint(500)) * time.Millisecond)
		ws := c.wireStat()
		conn.Write(ws.Bytes())
		// pretty.Log("STAT", ws.peerStat())
	}
}
