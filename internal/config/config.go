package config

import (
	"errors"
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		HTTP `yaml:"http"`
		PG   `yaml:"postgres"`
		Nuts `yaml:"nuts"`
	}

	HTTP struct {
		Port string `yaml:"port"`
	}

	PG struct {
		DSN                string `yaml:"dsn" env:"PG_DSN"`
		MaxOpenConnections int32  `yaml:"max_connections"  env:"PG_MAX_CONNECT"`
	}
	Nuts struct {
		ClusterID   string `yaml:"cluster_id" env:"CLUSTER_ID"`
		ClientSubID string `yaml:"client_sub_id" env:"CLIENT_SUB_ID"`
		ClientPubId string `yaml:"client_pub_id" env:"CLIENT_PUB_ID"`

		Subject string `yaml:"subject" env:"SUBJECT"`
		URL     string `yaml:"url" env:"URL"`
	}
)

func NewConfig() (*Config, error) {
	path := fetchConfigPath()
	if path == "" {
		return nil, errors.New("config path is empty")
	}

	cfg := &Config{}
	err := cleanenv.ReadConfig(path, cfg)
	if err != nil {
		return nil, err
	}

	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func fetchConfigPath() string {
	var res string
	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()
	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}

func (c *Config) GetDBConfig() (*pgxpool.Config, error) {
	poolConfig, err := pgxpool.ParseConfig(c.PG.DSN)
	if err != nil {
		return nil, err
	}

	poolConfig.ConnConfig.BuildStatementCache = nil
	poolConfig.ConnConfig.PreferSimpleProtocol = true
	poolConfig.MaxConns = c.PG.MaxOpenConnections

	return poolConfig, nil
}
