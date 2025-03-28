package grpc

import (
	"net"

	fileservice "github.com/TOMMy-Net/tages/internal/grpc/file_service"
	"google.golang.org/grpc"
)

type Server struct {
	Options []grpc.ServerOption
	Server  *grpc.Server
}

func NewServer(gs ...grpc.ServerOption) *Server {
	server := grpc.NewServer(gs...)
	return &Server{
		Options: gs,
		Server:  server,
	}
}

func (s *Server) Serve(network string, port string) error {
	err := fileservice.Register(s.Server, fileservice.NewFileServer("./files"))
	if err != nil {
		return err
	}
	l, err := net.Listen(network, port)
	if err != nil {
		return err
	}

	if err := s.Server.Serve(l); err != nil {
		return err
	}
	return nil
}
