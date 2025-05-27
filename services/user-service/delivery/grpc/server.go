package grpc

import (
	"context"
	"fmt"
	"user-service/internal/usecase"
	"proto/userpb"
)

type UserServer struct {
	userpb.UnimplementedUserServiceServer
	UserUsecase *usecase.UserUsecase
}

func NewUserServer(userUsecase *usecase.UserUsecase) *UserServer {
	return &UserServer{UserUsecase: userUsecase}
}

func (s *UserServer) RegisterUser(ctx context.Context, req *userpb.UserRequest) (*userpb.UserResponse, error) {
	userID, err := s.UserUsecase.RegisterUser(req.Username, req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to register user: %v", err)
	}
	return &userpb.UserResponse{UserId: userID, Success: true}, nil
}

func (s *UserServer) AuthenticateUser(ctx context.Context, req *userpb.AuthRequest) (*userpb.UserResponse, error) {
	userID, err := s.UserUsecase.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %v", err)
	}
	return &userpb.UserResponse{UserId: userID, Success: true}, nil
}

func (s *UserServer) GetUserProfile(ctx context.Context, req *userpb.UserID) (*userpb.UserProfile, error) {
	user, err := s.UserUsecase.GetUserProfile(req.UserId)
	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}
	return &userpb.UserProfile{
		UserId:   user.UserID,
		Username: user.Username,
	}, nil
}
