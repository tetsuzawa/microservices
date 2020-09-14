package post

import (
	"context"
	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	"github.com/tetsuzawa/microservices/backend/pkg/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sync"
)

// MockDB - テスト・開発用のDB
type MockDB struct {
	mu    sync.RWMutex
	posts *MockPostsTable
}

// MockGateway - MockDBのアダプターの構造体
type MockGateway struct {
	db *MockDB
}

// NewMockGateway - MockDBのアダプターの構造体のコンストラクタ
func NewMockGateway(db *MockDB) Repository {
	return &MockGateway{db}
}

// MockPostTable - ポストテーブル
type MockPostsTable struct {
	data map[string]api.Post
}

func newMockPostsTable() *MockPostsTable {
	return &MockPostsTable{data: make(map[string]api.Post)}
}

// CreatePost - Postを登録
func (r *MockGateway) CreatePost(ctx context.Context, userID, text string) (api.Post, error) {
	r.db.mu.Lock()
	defer r.db.mu.Unlock()

	// generate ID
	u, err := uuid.NewRandom()
	if err != nil {
		return api.Post{}, status.Error(codes.Unknown, "failed to insert into post-> "+err.Error())
	}
	id := u.String()

	post := api.Post{
		Id:           id,
		UserId:       userID,
		Text:         text,
		ParentPostId: "",
		CommentCount: 0,
		CreatedAt:    ptypes.TimestampNow(),
		UpdatedAt:    ptypes.TimestampNow(),
	}

	r.db.posts.data[id] = post

	return post, nil
}

// GetPostByID - 指定したIDのPostを取得
func (r *MockGateway) GetPostByID(ctx context.Context, id string) (api.Post, error) {
	post, ok := r.db.posts.data[id]
	if !ok {
		return api.Post{}, status.Error(codes.Unknown, "failed to read Post")
	}
	return post, nil
}

// UpdatePost - Postを更新
func (r *MockGateway) UpdatePost(ctx context.Context, id, userID, text string) (api.Post, error) {
	r.db.mu.Lock()
	defer r.db.mu.Unlock()

	post := r.db.posts.data[id]
	if post.UserId != userID {
		return api.Post{}, status.Error(codes.InvalidArgument, "User ID is not valid")
	}
	post.Text = text
	post.UpdatedAt = ptypes.TimestampNow()
	r.db.posts.data[post.Id] = post
	return post, nil
}

// DeletePost - Postを削除
func (r *MockGateway) DeletePost(ctx context.Context, id, userID string) (bool, error) {
	r.db.mu.Lock()
	defer r.db.mu.Unlock()

	post := r.db.posts.data[id]
	if post.UserId != userID {
		return false, status.Error(codes.InvalidArgument, "User ID is not valid")
	}
	delete(r.db.posts.data, id)
	return true, nil
}

// ListPosts - Postをすべて取得
func (r *MockGateway) ListPosts(ctx context.Context) ([]*api.Post, error) {
	r.db.mu.RLock()
	defer r.db.mu.RUnlock()

	var posts []*api.Post
	for _, v := range r.db.posts.data {
		posts = append(posts, &v)
	}
	return posts, nil
}
