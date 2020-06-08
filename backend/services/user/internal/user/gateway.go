package user

import (
	"database/sql"
	"github.com/gomodule/redigo/redis"
)

type Gateway struct {
	db        *sql.DB
	SessionDB redis.Conn
}

func NewGateway(db *sql.DB, sessionDB redis.Conn) Repository {
	return &Gateway{db, sessionDB}
}
