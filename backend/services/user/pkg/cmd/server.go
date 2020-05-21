package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	"github.com/kelseyhightower/envconfig"

	"github.com/tetsuzawa/microservices/backend/pkg/mysqlx"
	"github.com/tetsuzawa/microservices/backend/pkg/redisx"
	"github.com/tetsuzawa/microservices/backend/services/user/internal/user"
	"github.com/tetsuzawa/microservices/backend/services/user/protocol/grpc"
)

// Config - サーバのホストとポートのコンフィグ
type Config struct {
	// gRPC server start parameters section
	// gRPC is TCP port to listen by gRPC server
	GRPCHost string `split_words:"true"`
	GRPCPort string `split_words:"true"`
}

// ReadConfigEnv - APIサーバに関する設定を読み込む
func ReadConfigEnv(cfg *Config) error {
	err := envconfig.Process("API", cfg)
	if err != nil {
		return fmt.Errorf("envconfig.Process: %w", err)
	}
	if cfg.GRPCHost == "" {
		cfg.GRPCHost = "127.0.0.1"
	}
	if cfg.GRPCPort == "" {
		cfg.GRPCPort = "8080"
	}
	return nil
}

// RunServer runs gRPC server and HTTP gateway
func RunServer() error {
	ctx := context.Background()

	// get configuration
	var cfg Config
	err := ReadConfigEnv(&cfg)
	if err != nil {
		return fmt.Errorf("config read error: %w", err)
	}

	db, err := newDB()
	if err != nil {
		return err
	}
	defer db.Close()
	sessionDB, err := newSessionDB()
	if err != nil {
		return err
	}
	defer db.Close()
	repository := user.NewGateway(db, sessionDB)
	userServiceServer := user.NewUserServiceServer(repository)

	return grpc.RunServer(ctx, userServiceServer, cfg.GRPCHost, cfg.GRPCPort)
}

func newDB() (*gorm.DB, error) {
	// Mysql
	var cfg mysqlx.Config
	err := mysqlx.ReadEnv(&cfg)
	if err != nil {
		return nil, fmt.Errorf("mysqlx env read error %w", err)
	}
	log.Printf("Connecting to MySQL ...")
	db, err := cfg.Connect()
	if err != nil {
		return nil, fmt.Errorf("mysqlx connection error %w", err)
	}
	return db, nil
}

func newSessionDB() (redis.Conn, error) {
	// Redis
	var cfg redisx.Config
	err := redisx.ReadEnv(&cfg)
	if err != nil {
		return nil, fmt.Errorf("mysqlx env read error %w", err)
	}
	log.Printf("Connecting to REDIS ...")
	sessionDB, err := cfg.Connect()
	if err != nil {
		return nil, fmt.Errorf("redisx connection error %w", err)
	}
	return sessionDB, nil
}
