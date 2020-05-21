package user

import (
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
)

type Gateway struct {
	db        *gorm.DB
	SessionDB redis.Conn
}

func NewGateway(db *gorm.DB, sessionDB redis.Conn) Repository {
	return &Gateway{db, sessionDB}
}
