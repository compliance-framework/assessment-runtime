package plugin

import (
	"context"
	"log"
	"net"

	"github.com/compliance-framework/assessment-runtime/plugin/proto"

	"google.golang.org/grpc"
)

type ExecuteAction func(in *proto.ActionInput) (*proto.ActionOutput, error)

type pluginServer struct {
	proto.UnimplementedActionServiceServer
	ExecuteFunction func(in *proto.ActionInput) (*proto.ActionOutput, error)
}

func (s *pluginServer) Execute(ctx context.Context, in *proto.ActionInput) (*proto.ActionOutput, error) {
	if s.ExecuteFunction == nil {
		log.Output(2, "ExecuteFunction is nil")
		return nil, nil
	}
	return s.ExecuteFunction(in)
}

func Activate(executeAction ExecuteAction) {
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	proto.RegisterActionServiceServer(s, &pluginServer{ExecuteFunction: executeAction})
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
