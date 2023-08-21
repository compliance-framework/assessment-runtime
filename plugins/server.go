package plugins

import (
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

type Server struct {
	GrpcServer *grpc.Server
	Listener   net.Listener
	Plugin     Plugin
}

func (s *Server) Start() {
	go func() {
		if err := s.GrpcServer.Serve(s.Listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
}

func (s *Server) Stop() {
	s.GrpcServer.GracefulStop()
}
