package plugin

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/compliance-framework/assessment-runtime/plugin/proto"
	"google.golang.org/grpc"
)

type pluginServer struct {
	proto.UnimplementedActionServiceServer
	ExecuteFunc func(in *proto.ActionInput) (*proto.ActionOutput, error)
}

func NewPluginServer(executeFunc func(in *proto.ActionInput) (*proto.ActionOutput, error)) *pluginServer {
	return &pluginServer{ExecuteFunc: executeFunc}
}

func (s *pluginServer) Execute(ctx context.Context, in *proto.ActionInput) (*proto.ActionOutput, error) {
	if s.ExecuteFunc == nil {
		return nil, fmt.Errorf("ExecuteFunc is nil")
	}
	return s.ExecuteFunc(in)
}

func Activate(executeFunc func(in *proto.ActionInput) (*proto.ActionOutput, error)) {
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	proto.RegisterActionServiceServer(s, NewPluginServer(executeFunc))
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
