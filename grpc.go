//go:generate protoc -I . --go_out=plugins=grpc:./proto ./proto/hyparview.proto ./proto/gossip.proto

package main

type server struct{}
