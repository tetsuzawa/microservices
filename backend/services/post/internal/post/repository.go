package post

import (
	"context"
	"github.com/tetsuzawa/microservices/backend/pkg/api"
)

type Repository interface {
	Create(ctx context.Context, userID, text string) (api.Post, error)
}
