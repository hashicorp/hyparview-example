package main

import (
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Example is the main entry point
func main() {
	c, _ := MakeClient(&ClientConfig{
		ID:          makeID(),
		Addr:        os.Getenv("LISTEN"),
		Boot:        os.Getenv("BOOTSTRAP"),
		RootPEMFile: os.Getenv("ROOT_PEM"),
	})

	go c.lpShuffle()
	go c.joinCluster()
	server(c)
}

func server(c *Client) {
	lis, err := net.Listen("tcp", c.config.Addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	creds, err := credentials.NewServerTLSFromFile(c.config.RootPEMFile, c.config.KeyFile)
	if err != nil {
		log.Fatalf("Failed to generate credentials %v", err)
	}

	opts = []grpc.ServerOption{grpc.Creds(creds)}
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterRouteGuideServer(grpcServer, newServer(c))
	grpcServer.Serve(lis)
}
