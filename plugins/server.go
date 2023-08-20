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
	DoneCh     chan struct{}
}

func NewServer(plugin Plugin) *Server {
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	RegisterActionServiceServer(grpcServer, &ActionService{Plugin: plugin})

	return &Server{
		GrpcServer: grpcServer,
		Listener:   listener,
		Plugin:     plugin,
		DoneCh:     make(chan struct{}),
	}
}

func (s *Server) Start() {
	go func() {
		log.Tracef("Starting plugin server")
		if err := s.GrpcServer.Serve(s.Listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
}

func (s *Server) Stop() {
	log.Tracef("Stopping plugin server")
	s.GrpcServer.GracefulStop()
	log.Tracef("Stopped plugin server")
}
