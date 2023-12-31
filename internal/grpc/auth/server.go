package auth

import (
	"context"

	ssov1 "github.com/Alextek777/Bank/protos/gen/sso/sso"
	"google.golang.org/grpc"
)

type ServerAPI struct {
	ssov1.UnimplementedAuthServer
}

func Register(gRPC *grpc.Server) {
	ssov1.RegisterAuthServer(gRPC, &ServerAPI{})
}

func (s *ServerAPI) Login(
	ctx context.Context,
	req *ssov1.LoginRequest,
) (*ssov1.LoginResponse, error) {
	panic("implement me")
}

func (s *ServerAPI) Register(
	ctx context.Context,
	req *ssov1.RegisterRequest,
) (*ssov1.RegisterResponse, error) {
	panic("implement me")
}
