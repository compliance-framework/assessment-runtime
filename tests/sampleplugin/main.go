package main

import (
	"log"
	"net"

	"github.com/compliance-framework/assessment-runtime/plugin"
	"github.com/compliance-framework/assessment-runtime/plugin/proto"

	"google.golang.org/grpc"
)

func main() {
	println("gRPC server tutorial in Go")

	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	proto.RegisterActionServiceServer(s, &plugin.PluginServer{})
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
