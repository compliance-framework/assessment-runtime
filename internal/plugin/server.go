package plugin

import (
	"github.com/compliance-framework/assessment-runtime/plugin"
	"google.golang.org/grpc"
	"net"
	"os"

	log "github.com/sirupsen/logrus"
)

type Server struct {
	grpcServer   *grpc.Server
	listener     net.Listener
	plugin       plugin.Plugin
	shutdownChan chan struct{}
}

func NewServer(plugin plugin.Plugin) *Server {
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	RegisterActionServiceServer(grpcServer, &ActionService{plugin: plugin})

	return &Server{
		grpcServer:   grpcServer,
		listener:     listener,
		plugin:       plugin,
		shutdownChan: make(chan struct{}),
	}
}

//func Activate(executeFunc func(ctx context.Context, in *ActionInput) (*ActionOutput, error)) {
//
//	listener, err := net.Listen("tcp", ":9000")
//	if err != nil {
//		log.Fatalf("Failed to listen: %v", err)
//		return
//	}
//
//	log.Tracef("Starting plugin server")
//
//	defer func() {
//		_ = listener.Close()
//	}()
//
//	s := grpc.NewServer()
//	RegisterActionServiceServer(s, &ActionService{ExecuteFunc: executeFunc})
//
//	go func() {
//		if err := s.Serve(listener); err != nil {
//			log.Fatalf("failed to serve: %v", err)
//		}
//	}()
//
//	log.Tracef("Started plugin server")
//}

func Register(plugin plugin.Plugin) {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.TraceLevel)

}
