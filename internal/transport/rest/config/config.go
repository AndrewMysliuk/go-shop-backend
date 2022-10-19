package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/viper"
)

type Config struct {
	DB                DB
	Postgres          Postgres
	FileStorageConfig FileStorageConfig

	Server struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"server"`
}

type FileStorageConfig struct {
	Endpoint  string `envconfig:"STORAGE_ENDPOINT"`
	Bucket    string `envconfig:"STORAGE_BUCKET"`
	AccessKey string `envconfig:"STORAGE_ACCESS_KEY"`
	SecretKey string `envconfig:"STORAGE_SECRET_KEY"`
}

type DB struct {
	Host    string `envconfig:"DB_HOST"`
	Port    int    `envconfig:"DB_PORT"`
	SSLMode string `envconfig:"DB_SSLMODE"`
}

type Postgres struct {
	User     string `envconfig:"POSTGRES_USER"`
	Db       string `envconfig:"POSTGRES_DB"`
	Password string `envconfig:"POSTGRES_PASSWORD"`
}

func New(folder, filename string) (*Config, error) {
	cfg := new(Config)

	viper.AddConfigPath(folder)
	viper.SetConfigName(filename)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	if err := envconfig.Process("db", &cfg.DB); err != nil {
		return nil, err
	}

	if err := envconfig.Process("postgres", &cfg.Postgres); err != nil {
		return nil, err
	}

	if err := envconfig.Process("storage", &cfg.FileStorageConfig); err != nil {
		return nil, err
	}

	return cfg, nil
}
