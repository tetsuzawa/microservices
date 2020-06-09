package post

import (
	"context"
	"github.com/tetsuzawa/microservices/backend/pkg/api"
)

type PostServiceServer struct {
	r Repository
}

func NewPostServiceServer(r Repository) api.PostServiceServer {
	return &PostServiceServer{r}
}

func (s *PostServiceServer) CreatePost(ctx context.Context, request *api.CreatePostRequest) (*api.CreatePostResponse, error) {
	//TODO ユーザーの存在確認（別サービスとの通信）
	post, err := s.r.CreatePost(ctx, request.UserId, request.Text)
	if err != nil {
		return nil, err
	}
	return &api.CreatePostResponse{Post: &post}, nil
}

func (s *PostServiceServer) GetPost(ctx context.Context, request *api.GetPostRequest) (*api.GetPostResponse, error) {
	post, err := s.r.GetPostByID(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	return &api.GetPostResponse{Post: &post}, nil
}

func (s *PostServiceServer) UpdatePost(ctx context.Context, request *api.UpdatePostRequest) (*api.UpdatePostResponse, error) {
	//TODO ユーザーの存在確認（別サービスとの通信）
	post, err := s.r.UpdatePost(ctx, request.Id, request.UserId, request.Text)
	if err != nil {
		return nil, err
	}
	return &api.UpdatePostResponse{Post: &post}, nil
}

func (s *PostServiceServer) DeletePost(ctx context.Context, request *api.DeletePostRequest) (*api.DeletePostResponse, error) {
	//TODO ユーザーの存在確認（別サービスとの通信）
	isSuccess, err := s.r.DeletePost(ctx, request.Id, request.UserId)
	if err != nil {
		return nil, err
	}
	return &api.DeletePostResponse{IsSuccess: isSuccess}, nil
}

func (s *PostServiceServer) ListPosts(ctx context.Context, request *api.ListPostsRequest) (*api.ListPostsResponse, error) {
	posts, err := s.r.ListPosts(ctx)
	if err != nil {
		return nil, err
	}
	return &api.ListPostsResponse{Posts: posts}, nil
}
