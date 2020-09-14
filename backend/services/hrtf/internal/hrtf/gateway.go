package hrtf

import (
	"database/sql"

	"github.com/tetsuzawa/microservices/backend/pkg/awsx"
)

type Gateway struct {
	db      *sql.DB
	storage *awsx.Connection
}

func NewGateway(db *sql.DB, strg *awsx.Connection) Repository {
	return &Gateway{db, strg}
}
