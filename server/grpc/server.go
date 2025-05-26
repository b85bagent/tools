package server

import (
	"context"
	"fmt"
	"net"

	"github.com/b85bagent/tools/server"
	grpc "google.golang.org/grpc"
)

type GRPCServer struct {
	grpcServer *grpc.Server
	port       int
}

func NewGRPCServer(port int, grpcServer *grpc.Server) server.Server {
	return &GRPCServer{
		grpcServer: grpcServer,
		port:       port,
	}
}

func (s *GRPCServer) Start() error {
	addr := fmt.Sprintf(":%d", s.port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}
	return s.grpcServer.Serve(listener)
}

func (s *GRPCServer) Stop(ctx context.Context) error {
	done := make(chan struct{})
	go func() {
		s.grpcServer.GracefulStop()
		close(done)
	}()

	select {
	case <-ctx.Done():
		s.grpcServer.Stop()
		return ctx.Err()
	case <-done:
		return nil
	}
}
