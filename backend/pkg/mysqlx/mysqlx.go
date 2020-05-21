package mysqlx

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/kelseyhightower/envconfig"
)

// Config - Mysqlの接続情報に関する設定
type Config struct {
	User      string `split_words:"true"`
	Password  string `split_words:"true"`
	Protocol  string `split_words:"true"`
	Host      string `split_words:"true"`
	Port      string `split_words:"true"`
	DBName    string `split_words:"true"`
	Charset   string `split_words:"true"`
	ParseTime string `split_words:"true"`
	Loc       string `split_words:"true"`
}

// ReadEnv - 指定したenvfileからMysqlに関する設定を読み込む
func ReadEnv(cfg *Config) error {
	err := envconfig.Process("MYSQL", cfg)
	if err != nil {
		return fmt.Errorf("envconfig.Process: %w", err)
	}
	return nil
}

func (c *Config) build() {
	if c.User == "" {
		c.User = "root"
	}
	if c.Password == "" {
		c.Password = "root"
	}
	if c.Protocol == "" {
		c.Protocol = "tcp"
	}
	if c.Host == "" {
		c.Host = "127.0.0.1"
	}
	if c.Port == "" {
		c.Port = "3306"
	}
	if c.Charset == "" {
		c.Charset = "utf8mb4"
	}
	if c.ParseTime == "" {
		c.ParseTime = "True"
	}
	if c.Loc == "" {
		c.Loc = "Local"
	}
}

// Connect - Mysqlに接続
func (c *Config) Connect() (*gorm.DB, error) {
	c.build()
	const DBMS = "mysql"
	CONNECT := fmt.Sprintf(
		"%s:%s@%s(%s:%s)/%s?charset=%s&parseTime=%s&loc=%s",
		c.User,
		c.Password,
		c.Protocol,
		c.Host,
		c.Port,
		c.DBName,
		c.Charset,
		c.ParseTime,
		c.Loc,
	)

	db, err := gorm.Open(DBMS, CONNECT)
	if err != nil {
		return nil, fmt.Errorf("gorm.Open: %w", err)
	}
	return db, nil
}
