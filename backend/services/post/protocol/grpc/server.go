package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"

	"github.com/tetsuzawa/microservices/backend/pkg/api"
)

// RunServer runs gRPC service to publish service
func RunServer(ctx context.Context, postSrvc api.PostServiceServer, host string, port string) error {
	address := fmt.Sprintf("%s:%s", host, port)
	listen, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	// register service
	server := grpc.NewServer()
	api.RegisterPostServiceServer(server, postSrvc)

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			log.Println("shutting down gRPC server...")

			server.GracefulStop()

			<-ctx.Done()
		}
	}()

	// start gRPC server
	log.Printf("starting gRPC server at %s:%s...", host, port)
	return server.Serve(listen)
}
