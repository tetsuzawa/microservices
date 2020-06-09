package post

import (
	"database/sql"
)

type Gateway struct {
	db *sql.DB
}

func NewGateway(db *sql.DB) Repository {
	return &Gateway{db}
}
