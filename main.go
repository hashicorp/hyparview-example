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
	"google.golang.org/grpc/credentials"
)

// Example is the main entry point
func main() {
	addr := os.Getenv("ADDR")
	boot := os.Getenv("BOOTSTRAP")
	c, _ := newClient(&clientConfig{
		id:        newID(),
		addr:      addr,
		bootstrap: boot,
		serverPEM: os.Getenv("SERVER_PEM"),
		serverKey: os.Getenv("SERVER_KEY"),
		clientPEM: os.Getenv("CLIENT_PEM"),
		clientKey: os.Getenv("CLIENT_KEY"),
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
		log.Fatalf("failed to listen: %v", err)
	}

	creds, err := credentials.NewServerTLSFromFile(c.config.serverPEM, c.config.serverKey)
	if err != nil {
		log.Fatalf("Failed to generate credentials %v", err)
	}
	opts := []grpc.ServerOption{grpc.Creds(creds)}

	grpcServer := grpc.NewServer(opts...)
	srv := newServer(c)
	proto.RegisterHyparviewServer(grpcServer, srv)
	proto.RegisterGossipServer(grpcServer, srv)
	grpcServer.Serve(lis)
}
