package config

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string     `yaml:"env" env-default:"local"`
	HttpServer HTTPServer `yaml:"http_server"`
	Tarantool  Tarantool  `yaml:"tarantool"`
}

type HTTPServer struct {
	Address        string        `yaml:"address" env:"YAKVS_HTTP_HOST" env-default:"localhost:8080"`
	ReadTimeout    time.Duration `yaml:"read_timeout" env:"YAKVS_HTTP_READ_TIMEOUT"  env-default:"60s"`
	IdleTimeout    time.Duration `yaml:"idle_timeout" env:"YAKVS_HTTP_IDLE_TIMEOUT" env-default:"60s"`
	ContextTimeout time.Duration `yaml:"context_timeout" env:"YAKVS_HTTP_CONTEXT_TIMEOUT" env-default:"10s"`
	WriteTimeout   time.Duration `yaml:"write_timeout" env:"YAKVS_HTTP_WRITE_TIMEOUT" env-default:"60s"`
}

type Tarantool struct {
	Host     string        `yaml:"host" env:"YAKVS_TARANTOOL_HOST"`
	User     string        `yaml:"user" env:"YAKVS_TARANTOOL_USER"`
	Password string        `yaml:"password" env:"YAKVS_TARANTOOL_PASSWORD"`
	Timeout  time.Duration `yaml:"timeout" env:"YAKVS_TARANTOOL_TIMEOUT"`
}

// ParseConfig parses config from yaml to Config
func ParseConfig() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = ".env"
		log.Println("CONFIG_PATH is empty. parsing ENV")
	}

	if _, err := os.Stat(configPath); errors.Is(err, fs.ErrNotExist) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	cfg := &Config{}

	if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return cfg
}
