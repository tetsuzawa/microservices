package hrtf

import (
	"github.com/jinzhu/gorm"

	"github.com/tetsuzawa/microservices/backend/pkg/awsx"
)

type Gateway struct {
	db      *gorm.DB
	storage *awsx.Connection
}

func NewGateway(db *gorm.DB, strg *awsx.Connection) Repository {
	return &Gateway{db, strg}
}
