package usersrv

import (
	"context"
	protos "goim-pro/api/protos"
)

type userServer struct {}

func New() protos.UserServiceServer {
	return &userServer{}
}

func (s *userServer) Register(ctx context.Context, req *protos.User) (*protos.ServerCommonResponse, error) {
	panic("implement me")
}

func (s *userServer) Login(ctx context.Context, req *protos.User) (*protos.User, error) {
	panic("implement me")
}
