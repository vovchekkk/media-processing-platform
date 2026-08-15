package config

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env                 string `yaml:"env" env-default:"development"`
	StoragePath         string `yaml:"storage_path" env-required:"true"`
	HTTPServer          `yaml:"http_server"`
	TaskProcessorConfig `yaml:"task_processor"`
	DatabaseConfig      `yaml:"db"`
}

type HTTPServer struct {
	Host        string        `yaml:"host" env-default:"localhost"`
	Port        int           `yaml:"port" env-default:"8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func (s HTTPServer) Address() string {
	return net.JoinHostPort(s.Host, strconv.Itoa(s.Port))
}

type TaskProcessorConfig struct {
	ProcessingDuration time.Duration `yaml:"processing_duration" env-default:"30s"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host" env:"DB_HOST" env-required:"true"`
	User     string `yaml:"user" env:"DB_USER" env-required:"true"`
	Password string `yaml:"password" env:"DB_PASSWORD" env-required:"true"`
	DBName   string `yaml:"dbname" env:"DB_NAME" env-required:"true"`
	Port     string `yaml:"port" env:"DB_PORT" env-required:"true"`
	SSLMode  string `yaml:"sslmode" end:"DB_SSL_MODE" env-required:"true"`
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
