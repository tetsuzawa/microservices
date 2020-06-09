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
	post, err := s.r.Create(ctx, request.UserId, request.Text)
	if err != nil {
		return nil, err
	}
	return &api.CreatePostResponse{Post: &post}, nil
}

func (s *PostServiceServer) GetPost(ctx context.Context, request *api.GetPostRequest) (*api.GetPostResponse, error) {
	//TODO
	return &api.GetPostResponse{}, nil
}

func (s *PostServiceServer) UpdatePost(ctx context.Context, request *api.UpdatePostRequest) (*api.UpdatePostResponse, error) {
	//TODO
	return &api.UpdatePostResponse{}, nil
}

func (s *PostServiceServer) DeletePost(ctx context.Context, request *api.DeletePostRequest) (*api.DeletePostResponse, error) {
	//TODO
	return &api.DeletePostResponse{}, nil
}

func (s *PostServiceServer) ListPosts(ctx context.Context, request *api.ListPostsRequest) (*api.ListPostsResponse, error) {
	//TODO
	return &api.ListPostsResponse{}, nil
}
