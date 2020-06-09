package post

import (
	"context"
	"github.com/tetsuzawa/microservices/backend/pkg/api"
)

type Repository interface {
	CreatePost(ctx context.Context, userID, text string) (api.Post, error)
	GetPostByID(ctx context.Context, id string) (api.Post, error)
}
