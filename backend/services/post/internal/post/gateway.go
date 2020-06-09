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

// コネクションプールからSQL connectionを取得
func (r *Gateway) connect(ctx context.Context) (*sql.Conn, error) {
	c, err := r.db.Conn(ctx)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to connect to database-> "+err.Error())
	}
	return c, nil
}

// CreatePost - Postを作成
func (r *Gateway) CreatePost(ctx context.Context, userID, text string) (api.Post, error) {
	// SQL connectionを取得
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

	// Postを入れる変数を宣言
	var post api.Post
	// Transactions
	trans := func(tx *sql.Tx) error {
		// IDを生成
		u, err := uuid.NewRandom()
		if err != nil {
			return status.Error(codes.Unknown, "failed to insert into post-> "+err.Error())
		}
		id := u.String()

		// Postを挿入
		_, err = tx.ExecContext(ctx, "INSERT INTO posts (id, user_id, text) VALUES(?, ?, ?)",
			id, userID, text)
		if err != nil {
			return status.Error(codes.Unknown, "failed to insert into post-> "+err.Error())
		}

		// 作成したPostを取得
		row := tx.QueryRow("SELECT * FROM posts WHERE id = ?", id)
		err = row.Scan(&post.Id, &post.UserId, &post.Text, &post.CommentCount, &post.ParentPostId, &post.CreatedAt, &post.UpdatedAt)
		if err == sql.ErrNoRows {
			return status.Error(codes.InvalidArgument, "post not found-> "+err.Error())
		} else if err != nil {
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

// GetPostByID - 指定したIDのPostを取得
func (r *Gateway) GetPostByID(ctx context.Context, id string) (api.Post, error) {
	// SQL connectionを取得
	c, err := r.connect(ctx)
	if err != nil {
		return api.Post{}, err
	}
	defer c.Close()

	// Postを入れる変数を宣言
	var post api.Post
	// Post ID から Postを取得
	row := c.QueryRowContext(ctx, "SELECT * FROM posts WHERE id = ?", id)
	err = row.Scan(&post.Id, &post.UserId, &post.Text, &post.CommentCount, &post.ParentPostId, &post.CreatedAt, &post.UpdatedAt)
	if err == sql.ErrNoRows {
		return api.Post{}, status.Error(codes.InvalidArgument, "post not found-> "+err.Error())
	} else if err != nil {
		return api.Post{}, status.Error(codes.Unknown, "failed to read Post-> "+err.Error())
	}
	return post, nil
}

// UpdatePost - Postを更新
func (r *Gateway) UpdatePost(ctx context.Context, id, userID, text string) (api.Post, error) {
	// SQL connectionを取得
	c, err := r.connect(ctx)
	if err != nil {
		return api.Post{}, err
	}
	defer c.Close()

	//Transactionを開始
	tx, err := c.BeginTx(ctx, nil)
	if err != nil {
		if err != nil {
			return api.Post{}, status.Error(codes.Unknown, "failed to update Post-> "+err.Error())
		}
	}

	// Postを入れる変数を宣言
	var post api.Post
	// Transactions
	trans := func(tx *sql.Tx) error {

		// 更新前のPostを取得
		row := tx.QueryRow("SELECT * FROM posts WHERE id = ?", id)
		err = row.Scan(&post.Id, &post.UserId, &post.Text, &post.CommentCount, &post.ParentPostId, &post.CreatedAt, &post.UpdatedAt)
		if err == sql.ErrNoRows {
			return status.Error(codes.InvalidArgument, "post not found-> "+err.Error())
		} else if err != nil {
			return status.Error(codes.Unknown, "failed to update Post-> "+err.Error())
		}
		if post.UserId != userID {
			return status.Error(codes.InvalidArgument, "User ID is not valid")
		}

		// Postを更新
		_, err = tx.Exec("UPDATE posts SET text = ? WHERE id = ?", text, id)
		if err != nil {
			return status.Error(codes.Unknown, "failed to update Post-> "+err.Error())
		}

		// 更新後のPostを取得
		row = tx.QueryRow("SELECT * FROM posts WHERE id = ?", id)
		err = row.Scan(&post.Id, &post.UserId, &post.Text, &post.CommentCount, &post.ParentPostId, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return status.Error(codes.Unknown, "failed to read updated Post-> "+err.Error())
		}

		return nil
	}

	if err := trans(tx); err != nil {
		if re := tx.Rollback(); re != nil {
			err = fmt.Errorf("%w -> %s", err, re.Error())
		}
		return api.Post{}, status.Error(codes.Unknown, "failed to update Post-> "+err.Error())
	}

	if err = tx.Commit(); err != nil {
		return api.Post{}, status.Error(codes.Unknown, "failed to update Post-> "+err.Error())
	}
	return post, nil
}

// DeletePost - Postを削除
func (r *Gateway) DeletePost(ctx context.Context, id, userID string) (bool, error) {
	// SQL connectionを取得
	c, err := r.connect(ctx)
	if err != nil {
		return false, err
	}
	defer c.Close()

	//Transactionを開始
	tx, err := c.BeginTx(ctx, nil)
	if err != nil {
		return false, status.Error(codes.Unknown, "failed to delete Post-> "+err.Error())
	}

	// Transactions
	trans := func(tx *sql.Tx) error {
		// Postを入れる変数を宣言
		var post api.Post

		// 削除前のPostを取得
		row := tx.QueryRow("SELECT * FROM posts WHERE id = ?", id)
		err = row.Scan(&post.Id, &post.UserId, &post.Text, &post.CommentCount, &post.ParentPostId, &post.CreatedAt, &post.UpdatedAt)
		if err == sql.ErrNoRows {
			return status.Error(codes.InvalidArgument, "post not found-> "+err.Error())
		} else if err != nil {
			return status.Error(codes.Unknown, "failed to delete Post-> "+err.Error())
		}

		if post.UserId != userID {
			return status.Error(codes.InvalidArgument, "User ID is not valid")
		}

		// Postを削除
		_, err = tx.Exec("DELETE FROM posts WHERE id = ?", id)
		if err != nil {
			return status.Error(codes.Unknown, "failed to delete Post-> "+err.Error())
		}

		return nil
	}

	if err := trans(tx); err != nil {
		if re := tx.Rollback(); re != nil {
			err = fmt.Errorf("%w -> %s", err, re.Error())
		}
		return false, status.Error(codes.Unknown, "failed to update Post-> "+err.Error())
	}

	if err = tx.Commit(); err != nil {
		return false, status.Error(codes.Unknown, "failed to update Post-> "+err.Error())
	}
	return true, nil
}

// ListPosts - Postをすべて取得
func (r *Gateway) ListPosts(ctx context.Context) ([]*api.Post, error) {
	// SQL connectionを取得
	c, err := r.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	// Postを入れる変数を宣言
	// Post ID から Postを取得
	rows, err := c.QueryContext(ctx, "SELECT * FROM posts ORDER BY created_at")
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to list Posts-> "+err.Error())
	}

	var posts []*api.Post
	for rows.Next() {
		var post api.Post
		err = rows.Scan(&post.Id, &post.UserId, &post.Text, &post.CommentCount, &post.ParentPostId, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return nil, status.Error(codes.Unknown, "failed to list Posts-> "+err.Error())
		}
		posts = append(posts, &post)
	}

	return posts, nil
}
