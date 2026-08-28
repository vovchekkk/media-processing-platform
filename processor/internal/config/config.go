package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env            string `yaml:"env" env-default:"development"`
	StoragePath    string `yaml:"storage_path" env-required:"true"`
	DatabaseConfig `yaml:"db"`
	RabbitMQConfig `yaml:"rabbitmq"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host" env:"DB_HOST" env-required:"true"`
	User     string `yaml:"user" env:"DB_USER" env-required:"true"`
	Password string `yaml:"password" env:"DB_PASSWORD" env-required:"true"`
	DBName   string `yaml:"dbname" env:"DB_NAME" env-required:"true"`
	Port     string `yaml:"port" env:"DB_PORT" env-required:"true"`
	SSLMode  string `yaml:"sslmode" env:"DB_SSL_MODE" env-required:"true"`
}

type RabbitMQConfig struct {
	Host       string                   `yaml:"host" env:"RABBITMQ_HOST" env-required:"true"`
	Port       string                   `yaml:"port" env:"RABBITMQ_PORT" env-required:"true"`
	User       string                   `yaml:"user" env:"RABBITMQ_USER" env-required:"true"`
	Password   string                   `yaml:"password" env:"RABBITMQ_PASSWORD" env-required:"true"`
	QueueName  string                   `yaml:"queue_name" env:"RABBITMQ_QUEUE_NAME" env-required:"true"`
	Connection RabbitMQConnectionConfig `yaml:"connection"`
}

type RabbitMQConnectionConfig struct {
	InitialBackoff    time.Duration `yaml:"initial_backoff" env-default:"1s"`
	MaxRetries        int           `yaml:"max_retries" env-default:"5"`
	ReconnectInterval time.Duration `yaml:"reconnect_interval" env-default:"2s"`
}

func (dbConfig DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		dbConfig.Host,
		dbConfig.User,
		dbConfig.Password,
		dbConfig.DBName,
		dbConfig.Port,
		dbConfig.SSLMode)
}

func (r RabbitMQConfig) DSN() string {
	return fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		r.User,
		r.Password,
		r.Host,
		r.Port,
	)
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable is not set")
	}

	if _, err := os.Stat(configPath); err != nil {
		log.Fatalf("error opening config file: %s", err)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("error reading config file: %s", err)
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("error reading env variables: %s", err)
	}

	return &cfg
}
