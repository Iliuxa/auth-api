package authgrpc

import (
	"auth-api/internal/domain"
	"auth-api/internal/usecase"
	"context"
	"errors"
	proto "github.com/Iliuxa/protos/gen/proto"
	"github.com/go-playground/validator/v10"
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
	validate := *validator.New()
	var isValidateErr = false
	isValidateErr = validate.Var(in.GetEmail(), "email,min=1") != nil
	isValidateErr = isValidateErr || validate.Var(in.GetPassword(), "min=1") != nil
	if isValidateErr {
		return nil, status.Error(codes.InvalidArgument, "Validation Error")
	}

	token, err := s.auth.Login(ctx, in.GetEmail(), in.GetPassword())

	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) || errors.Is(err, domain.ErrUserNotFound) {
			return nil, status.Error(codes.InvalidArgument, "Invalid email or password")
		}

		return nil, status.Error(codes.Internal, "Failed to login")
	}

	return &proto.LoginResponse{Jwt: token}, nil
}

func (s *serverAPI) Register(ctx context.Context, in *proto.RegisterRequest) (*proto.LoginResponse, error) {
	validate := *validator.New()
	var isValidateErr = false
	isValidateErr = validate.Var(in.GetLogin().GetEmail(), "email,min=1") != nil
	isValidateErr = isValidateErr || validate.Var(in.GetLogin().GetPassword(), "min=1") != nil
	isValidateErr = isValidateErr || validate.Var(in.GetName(), "min=1") != nil
	if isValidateErr {
		return nil, status.Error(codes.InvalidArgument, "Validation Error")
	}

	token, err := s.auth.Register(ctx, in.GetLogin().GetEmail(), in.GetLogin().GetPassword(), in.GetName())

	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			return nil, status.Error(codes.InvalidArgument, "User already exists")
		}

		return nil, status.Error(codes.Internal, "Failed to register")
	}

	return &proto.LoginResponse{Jwt: token}, nil
}
