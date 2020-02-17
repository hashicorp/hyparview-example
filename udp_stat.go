package main

import (
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

type stats struct {
	safe map[string]*peerStat
	lock sync.RWMutex
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
			stats.lock.Lock()
			stats.safe[p.From] = p
			stats.lock.Unlock()
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
	as := strings.Split(addr, ":")
	ip := net.ParseIP(as[0])
	port, _ := strconv.Atoi(as[1])
	a := &net.UDPAddr{IP: ip, Port: port}
	conn, err := net.ListenUDP("udp", a)
	if err != nil {
		log.Fatal("error: udp listen %v", err)
	}
	conn.SetReadBuffer(c.config.statUDPBuffer)
	defer conn.Close()

	buf := make([]byte, STAT_SIZE)
	for {
		conn.ReadFromUDP(buf)
		inbox <- buf
	}
}

// runStatClient blocks forever sending stats on a random interval to the remote addr
func runStatClient(c *client, addr string) {
	log.Printf("debug: stat client %s", addr)

	// Under the impression that udp dial won't return an error :shrug:
	conn, err := net.Dial("udp", addr)
	for {
		if err == nil {
			break
		}
		time.Sleep(time.Second)
		conn, err = net.Dial("udp", addr)
	}
	defer conn.Close()

	for {
		c.statSleep()
		ws := c.wireStat()
		conn.Write(ws.Bytes())
	}
}
