package post

import (
	"context"
	"database/sql"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
	"github.com/tetsuzawa/microservices/backend/pkg/api"
)

type Gateway struct {
	db *sql.DB
}

func NewGateway(db *sql.DB) Repository {
	return &Gateway{db}
}

// connect returns SQL database connection from the pool
func (r *Gateway) connect(ctx context.Context) (*sql.Conn, error) {
	c, err := r.db.Conn(ctx)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to connect to database-> "+err.Error())
	}
	return c, nil
}

func (r *Gateway) Create(ctx context.Context, userID, text string) (api.Post, error) {
	// get SQL connection from pool
	c, err := r.connect(ctx)
	if err != nil {
		return api.Post{}, err
	}
	defer c.Close()

	tx, err := c.BeginTx(ctx, nil)
	if err != nil {
		if err != nil {
			return api.Post{}, status.Error(codes.Unknown, "failed to insert into post-> "+err.Error())
		}
	}

	// declare post to return if success
	var post api.Post
	// Transactions
	trans := func(tx *sql.Tx) error {
		// generate ID
		u, err := uuid.NewRandom()
		if err != nil {
			return status.Error(codes.Unknown, "failed to insert into post-> "+err.Error())
		}
		id := u.String()

		// insert Post entity data
		_, err = tx.ExecContext(ctx, "INSERT INTO posts (id, user_id, text) VALUES(?, ?, ?)",
			id, userID, text)
		if err != nil {
			return status.Error(codes.Unknown, "failed to insert into post-> "+err.Error())
		}

		// read created Post
		row := tx.QueryRow("SELECT * FROM posts WHERE id = ?", id)
		err = row.Scan(&post.Id, &post.UserId, &post.Text, &post.CommentCount, &post.ParentPostId, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return status.Error(codes.Unknown, "failed to read created Post-> "+err.Error())
		}
		return nil
	}

	if err := trans(tx); err != nil {
		if re := tx.Rollback(); re != nil {
			err = fmt.Errorf("%w -> %s", err, re.Error())
		}
		return api.Post{}, status.Error(codes.Unknown, "failed to read created Post-> "+err.Error())
	}

	if err = tx.Commit(); err != nil {
		return api.Post{}, status.Error(codes.Unknown, "failed to insert into post-> "+err.Error())
	}
	return post, nil
}
