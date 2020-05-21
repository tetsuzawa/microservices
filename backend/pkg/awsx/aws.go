package awsx

import (
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/kelseyhightower/envconfig"
)

// Config - AWSの接続情報に関する設定
type Config struct {
	AccessKey       string `split_words:"true"`
	SecretAccessKey string `split_words:"true"`
	S3Bucket        string `split_words:"true"`
	Region          string `split_words:"true"`
}

type Connection struct {
	Config  Config
	Session *session.Session
	SVC     *s3.S3
}

// ReadEnv - 指定したenvfileからAWSに関する設定を読み込む
func ReadEnv(cfg *Config) error {
	err := envconfig.Process("AWS", cfg)
	if err != nil {
		return fmt.Errorf("envconfig.Process: %w", err)
	}
	return nil
}

func (c *Config) build() {
	if c.AccessKey == "" {
		err := errors.New("error: env AWS_ACCESS_KEY is empty")
		log.Fatalln(err)
	}
	if c.SecretAccessKey == "" {
		err := errors.New("error: env AWS_SECRET_ACCESS_KEY is empty")
		log.Fatalln(err)
	}
	if c.Region == "" {
		c.Region = "ap-northeast-1"
	}
	if c.S3Bucket == "" {
		log.Fatalln("c.S3Bucket is empty:", c.S3Bucket)
	}
}

// Connect - AWSに接続
func (c *Config) Connect() (*Connection, error) {
	c.build()
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(c.AccessKey, c.SecretAccessKey, ""),
		Region:      aws.String(c.Region),
	})
	if err != nil {
		return nil, fmt.Errorf("session.NewSessionWithOptions: %w", err)
	}
	svc := s3.New(sess)
	return &Connection{Config: c, Session: sess, SVC: svc}, nil
}
