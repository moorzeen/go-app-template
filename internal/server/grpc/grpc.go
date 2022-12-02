package grpc

import (
	"google.golang.org/grpc"
	"net"
	"template/internal/service/auth"
	pbAuth "template/internal/service/auth/proto"
	"template/internal/service/user"
	pb "template/internal/service/user/proto"
)

type Server struct {
	server   *grpc.Server
	errCh    chan error
	listener net.Listener
}

func NewGrpcServer(userService *user.Service, authService *auth.Service, port string) (*Server, error) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return nil, err
	}
	server := grpc.NewServer()
	pb.RegisterUserServiceServer(server, userService)
	pbAuth.RegisterAuthServiceServer(server, authService)

	return &Server{
		server:   server,
		listener: lis,
		errCh:    make(chan error),
	}, nil
}

func (srv *Server) Start() {
	go func() {
		srv.errCh <- srv.server.Serve(srv.listener)
	}()
}

func (srv *Server) Stop() {
	srv.server.GracefulStop()
}

func (srv *Server) Error() chan error {
	return srv.errCh
}
