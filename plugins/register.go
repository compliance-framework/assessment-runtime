package plugins

import (
	"context"
	"google.golang.org/grpc"
	"net"
	"os"

	log "github.com/sirupsen/logrus"
)

func Register(plugin Plugin) {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.TraceLevel)

	listener, err := net.Listen("tcp", ":0")
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
	}

	log.WithFields(log.Fields{
		"network": server.Listener.Addr().Network(),
		"host":    listener.Addr().(*net.TCPAddr).IP.String(),
		"port":    listener.Addr().(*net.TCPAddr).Port,
	}).Info()

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
