package main

import (
	"log"
	"math"
	"math/rand"
	"net"
	"os"
	"time"

	h "github.com/hashicorp/hyparview"
	"github.com/hashicorp/hyparview-example/proto"
	"google.golang.org/grpc"
)

// Example is the main entry point
func main() {
	addr := os.Getenv("ADDR")
	boot := os.Getenv("BOOTSTRAP")
	http := os.Getenv("HTTP_UI")
	stat := os.Getenv("STAT_UDP")

	seed := h.Rint64Crypto(math.MaxInt64 - 1)
	rand.Seed(seed)

	c := newClient(&clientConfig{
		id:         newID(),
		addr:       addr,
		bootstrap:  boot,
		caCert:     os.Getenv("CA_CERT"),
		serverCert: os.Getenv("SERVER_CERT"),
		serverKey:  os.Getenv("SERVER_KEY"),
		clientCert: os.Getenv("CLIENT_CERT"),
		clientKey:  os.Getenv("CLIENT_KEY"),

		sendFanOut:     20,
		shuffleSeconds: 30,
		hvClientCount:  1000,
		hvInboxBuffer:  1024,
		hvOutboxBuffer: 1024,
		gossipMaxHeat:  5,

		statParseFanOut:  4,
		statParseBuffer:  100,
		statMillis:       5000,
		statUpdateBuffer: 1024,
		statUDPBuffer:    1048576,
	})

	go runServer(c)
	goRunClient(c)

	if boot != addr {
		for {
			ms := c.hv.SendJoin(node(boot))
			err := c.send(ms[0])
			if err == nil {
				break
			}
			time.Sleep(time.Second)
			log.Printf("debug: bootstrap fail %v, retrying", err)
		}
	}

	if http != "" {
		stats := newStats()
		go runStatServer(stat, stats, c)
		go runUIServer(http, c, stats)
	} else {
		go runStatClient(c, stat)
	}

	c.runShuffle()
}

// runShuffle blocks, sending a shuffle request to a random peer at some random point in the
// shuffleSeconds window
func (c *client) runShuffle() {
	for {
		s := time.Duration(h.Rint(c.config.shuffleSeconds))
		time.Sleep(s * time.Second)
		peer := c.getPeer()
		if peer == nil {
			continue
		}

		err := c.send(c.hv.SendShuffle(peer))
		if err != nil {
			c.failActive(peer)
			continue
		}
	}
}

// runServer blocks, starts grpc
func runServer(c *client) {
	lis, err := net.Listen("tcp", c.config.addr)
	if err != nil {
		log.Fatalf("error listen: %v\n", err)
	}

	creds, err := serverCreds(c.config)
	if err != nil {
		log.Fatalf("error tls: %v\n", err)
	}
	opts := []grpc.ServerOption{grpc.Creds(creds)}

	go runServerUpdater(c)

	grpcServer := grpc.NewServer(opts...)
	srv := newServer(c)
	proto.RegisterHyparviewServer(grpcServer, srv)
	proto.RegisterGossipServer(grpcServer, srv)
	grpcServer.Serve(lis)
}

// runServerUpdater blocks, apply hv updates
func runServerUpdater(c *client) {
	for {
		m := <-c.in
		c.recv(m)
	}
}

// goRunClient fans hv message sending processes
func goRunClient(c *client) {
	// fan out for sending
	for i := 0; i < c.config.sendFanOut; i++ {
		go func() {
			for {
				m := <-c.out
				c.send(m)
			}
		}()
	}
}
