package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/kelseyhightower/envconfig"

	"github.com/tetsuzawa/microservices/backend/bff/cmd/server/controller"
	"github.com/tetsuzawa/microservices/backend/pkg/middleware"
)

// APIConfig - APIサーバのホストとポートのコンフィグ
type APIConfig struct {
	Host string `split_words:"true"`
	Port string `split_words:"true"`
}

// ReadAPIEnv - APIサーバに関する設定を読み込む
func ReadAPIEnv(cfg *APIConfig) error {
	err := envconfig.Process("API", cfg)
	if err != nil {
		return fmt.Errorf("envconfig.Process: %w", err)
	}
	if cfg.Host == "" {
		cfg.Host = "127.0.0.1"
	}
	if cfg.Port == "" {
		cfg.Port = "8080"
	}
	return nil
}

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	// default
	var apiCfg APIConfig
	err := ReadAPIEnv(&apiCfg)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	h := &ServerHandler{RootHandler: &controller.RootHandler{}}
	router := middleware.RequestLog(h)
	router = middleware.ApiMakeHandler(router, "/api/")

	if err := RunServer(apiCfg.Host, apiCfg.Port, router); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func RunServer(host string, port string, handler http.Handler) error {
	//http.HandleFunc("/api/ping", controller.HandlePing)

	address := fmt.Sprintf("%s:%s", host, port)
	fmt.Printf("Listening on %s\n", address)
	return http.ListenAndServe(address, handler)
}
