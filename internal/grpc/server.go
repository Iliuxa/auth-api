package authgrpc

import (
	"auth-api/internal/domain"
	"auth-api/internal/usecase"
	"context"
	"errors"
	proto "github.com/Iliuxa/protos/gen/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	proto.UnimplementedAuthServiceServer
	auth usecase.AuthUsecase
}

func Register(gRPCServe *grpc.Server, auth usecase.AuthUsecase) {
	proto.RegisterAuthServiceServer(gRPCServe, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(ctx context.Context, in *proto.LoginInfo) (*proto.LoginResponse, error) {
	// todo validation

	token, err := s.auth.Login(ctx, in.GetEmail(), in.GetPassword())

	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "Invalid email or password")
		}

		return nil, status.Error(codes.Internal, "Failed to login")
	}

	return &proto.LoginResponse{Jwt: token}, nil
}

func (s *serverAPI) Register(ctx context.Context, in *proto.RegisterRequest) (*proto.LoginResponse, error) {
	// todo validation
	token, err := s.auth.Register(ctx, in.GetLogin().GetEmail(), in.GetLogin().GetPassword(), in.GetName())

	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			return nil, status.Error(codes.InvalidArgument, "User already exists")
		}

		return nil, status.Error(codes.Internal, "Failed to register")
	}

	return &proto.LoginResponse{Jwt: token}, nil
}
