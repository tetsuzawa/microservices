package grpc

import (
	"context"
	"github.com/tetsuzawa/microservices/backend/internal/user"
	"github.com/tetsuzawa/microservices/backend/pkg/api"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
)

func RunServer(ctx context.Context, srvs user.Services) error {
	listen, err := net.Listen("tcp", ":80")
	if err != nil {
		return err
	}

	// register service
	server := grpc.NewServer()
	api.RegisterUserServiceServer(server, srvs.UserServiceServer)
	api.RegisterHRTFServiceServer(server, srvs.HrtfServiceServer)

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
	log.Println("starting gRPC server...")
	return server.Serve(listen)
}