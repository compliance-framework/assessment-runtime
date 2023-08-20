package plugins

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"os"

	log "github.com/sirupsen/logrus"
)

func Register(plugin Plugin) {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.TraceLevel)

	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	RegisterActionServiceServer(grpcServer, &ActionService{Plugin: plugin})

	doneCh := make(chan struct{})

	server := &Server{
		GrpcServer: grpcServer,
		Listener:   listener,
		Plugin:     plugin,
		DoneCh:     doneCh,
	}

	fmt.Printf("%s|%s", server.Listener.Addr().Network(), server.Listener.Addr().String())

	go server.Start()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	select {
	case <-ctx.Done():
		_ = listener.Close()
		server.Stop()
		<-doneCh

	case <-doneCh:
		os.Exit(0)
	}
}
