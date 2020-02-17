package main

import (
	"log"
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
	c := newClient(&clientConfig{
		id:         newID(),
		addr:       addr,
		bootstrap:  boot,
		caCert:     os.Getenv("CA_CERT"),
		serverCert: os.Getenv("SERVER_CERT"),
		serverKey:  os.Getenv("SERVER_KEY"),
		clientCert: os.Getenv("CLIENT_CERT"),
		clientKey:  os.Getenv("CLIENT_KEY"),
	})

	go runServer(c)
	if boot != addr {
		for {
			err := c.send(h.SendJoin(node(boot), c.hv.Self))
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
		go runUIServer(http, stats)
	} else {
		go runStatClient(c, stat)
	}

	c.lpShuffle()
}

func (c *client) lpShuffle() {
	for {
		time.Sleep(10 * time.Second)
		peer := c.hv.Peer()
		if peer == nil {
			c.failActive(nil)
			continue
		}

		req := c.hv.SendShuffle(c.hv.Peer())
		r, err := c.sendShuffle(req)
		if err != nil {
			log.Printf("error: shuffle send: %v\n", err)
			continue
		}
		if r == nil {
			continue
		}

		c.hv.RecvShuffleReply(r)
	}
}

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

	grpcServer := grpc.NewServer(opts...)
	srv := newServer(c)
	proto.RegisterHyparviewServer(grpcServer, srv)
	proto.RegisterGossipServer(grpcServer, srv)
	grpcServer.Serve(lis)
}
