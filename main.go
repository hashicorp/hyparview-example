//go:generate protoc -I . --go_out=plugins=grpc:. ./proto/hyparview.proto ./proto/gossip.proto

package main

import (
	"log"
	"net"
	"os"
	"time"

	h "github.com/hashicorp/hyparview"
	"github.com/hashicorp/hyparview-example/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Example is the main entry point
func main() {
	boot := os.Getenv("BOOTSTRAP")
	c, _ := newClient(&clientConfig{
		ID:          newID(),
		Addr:        os.Getenv("LISTEN"),
		Boot:        boot,
		RootPEMFile: os.Getenv("ROOT_PEM"),
		keyFile:     os.Getenv("SERVER_KEY_PEM"),
	})

	go runServer(c)
	c.send(h.SendJoin(c.hv.Self, node(boot)))
	c.lpShuffle()
}

func (c *client) lpShuffle() {
	for {
		time.Sleep(10 * time.Second)
		r, err := c.sendShuffle(c.hv.SendShuffle(c.hv.Peer()))
		if err != nil {
			// log
			continue
		}
		c.hv.RecvShuffleReply(r)
	}
}

func runServer(c *client) {
	lis, err := net.Listen("tcp", c.config.Addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	creds, err := credentials.NewServerTLSFromFile(c.config.RootPEMFile, c.config.keyFile)
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
