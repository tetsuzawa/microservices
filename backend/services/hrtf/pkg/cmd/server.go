package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/kelseyhightower/envconfig"

	"github.com/tetsuzawa/microservices/backend/pkg/awsx"
	"github.com/tetsuzawa/microservices/backend/pkg/mysqlx"
	"github.com/tetsuzawa/microservices/backend/services/hrtf/internal/hrtf"
	"github.com/tetsuzawa/microservices/backend/services/hrtf/protocol/grpc"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
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
	strg, err := newStorage()
	if err != nil {
		return err
	}
	repository := hrtf.NewGateway(db, strg)
	hrtfServiceServer := hrtf.NewHRTFServiceServer(repository)

	return grpc.RunServer(ctx, hrtfServiceServer, cfg.GRPCHost, cfg.GRPCPort)
}

func newDB() (*sql.DB, error) {
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

func newStorage() (*awsx.Connection, error) {
	var cfg awsx.Config
	err := awsx.ReadEnv(&cfg)
	if err != nil {
		return nil, fmt.Errorf("awsx env read error %w", err)
	}
	log.Printf("Connecting to S3 ...")
	strg, err := cfg.Connect()
	if err != nil {
		return nil, fmt.Errorf("awsx connection error %w", err)
	}
	return strg, nil
}
