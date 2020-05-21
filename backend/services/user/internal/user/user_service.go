package user

import (
	"context"

	"github.com/tetsuzawa/microservices/backend/pkg/api"
)

type UserServiceServer struct {
	r Repository
}

func NewUserServiceServer(r Repository) api.UserServiceServer {
	return &UserServiceServer{r}
}

func (s *UserServiceServer) CreateUser(ctx context.Context, request *api.CreateUserRequest) (*api.CreateUserResponse, error) {
	//TODO
	return &api.CreateUserResponse{}, nil
}

func (s *UserServiceServer) GetUser(ctx context.Context, request *api.GetUserRequest) (*api.GetUserResponse, error) {
	//TODO
	return &api.GetUserResponse{}, nil
}

func (s *UserServiceServer) UpdateUser(ctx context.Context, request *api.UpdateUserRequest) (*api.UpdateUserResponse, error) {
	//TODO
	return &api.UpdateUserResponse{}, nil
}

func (s *UserServiceServer) DeleteUser(ctx context.Context, request *api.DeleteUserRequest) (*api.DeleteUserResponse, error) {
	//TODO
	return &api.DeleteUserResponse{}, nil
}
