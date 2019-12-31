package main

import (
	"log"
	"net"
	"os"

	"github.com/hashicorp/hyparview-example/proto/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Example is the main entry point
func main() {
	c, _ := newClient(&clientConfig{
		ID:          newID(),
		Addr:        os.Getenv("LISTEN"),
		Boot:        os.Getenv("BOOTSTRAP"),
		RootPEMFile: os.Getenv("ROOT_PEM"),
		keyFile:     os.Getenv("SERVER_KEY_PEM"),
	})

	go c.lpShuffle()
	go c.joinCluster()
	runServer(c)
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
