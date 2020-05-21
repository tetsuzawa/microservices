package hrtf

import "sync"

// MockDB - テスト・開発用のDB
type MockDB struct {
	mu     sync.RWMutex
}

type MockGateway struct {
	db         *MockDB
	storage    string
}
