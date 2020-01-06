//go:generate protoc -I . --go_out=plugins=grpc:. ./proto/hyparview.proto ./proto/gossip.proto

package main

import (
	"log"
	"net"
	"os"
	"time"

	h "github.com/hashicorp/hyparview"
	"github.com/hashicorp/hyparview-example/proto"
	"github.com/kr/pretty"
	"google.golang.org/grpc"
)

// Example is the main entry point
func main() {
	addr := os.Getenv("ADDR")
	boot := os.Getenv("BOOTSTRAP")
	c, _ := newClient(&clientConfig{
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
		c.send(h.SendJoin(node(boot), c.hv.Self))
	}
	c.lpShuffle()
}

func (c *client) lpShuffle() {
	for {
		time.Sleep(10 * time.Second)
		r, err := c.sendShuffle(c.hv.SendShuffle(c.hv.Peer()))
		if err != nil {
			pretty.Log("error lpShuffle", err)
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
		log.Fatalf("error listen: %v", err)
	}

	creds, err := serverCreds(c.config)
	if err != nil {
		log.Fatalf("error tls: %v", err)
	}
	opts := []grpc.ServerOption{grpc.Creds(creds)}

	grpcServer := grpc.NewServer(opts...)
	srv := newServer(c)
	proto.RegisterHyparviewServer(grpcServer, srv)
	proto.RegisterGossipServer(grpcServer, srv)
	grpcServer.Serve(lis)
}
