package app

import (
	saltyGrpc "goim-pro/internal/app/grpc"
	"google.golang.org/grpc"
)

type Server struct {
	grpcServer *grpc.Server
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) StartGrpcServer() *grpc.Server {
	return saltyGrpc.New()
}

func (s *Server) GracefulStopGrpcServer() {
	s.grpcServer.GracefulStop()
}

func (s *Server) ForceStopGrpcServer () {
	s.grpcServer.Stop()
}
