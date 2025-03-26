package config

import (
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
	Address        string        `yaml:"address" env-default:"localhost:8080"`
	ReadTimeout    time.Duration `yaml:"read_timeout" env-default:"4s"`
	IdleTimeout    time.Duration `yaml:"idle_timeout" env-default:"60s"`
	ContextTimeout time.Duration `yaml:"context_timeout" env-default:"10s"`
	WriteTimeout   time.Duration `yaml:"write_timeout" env-default:"60s"`
}

type Tarantool struct {
	Host     string        `yaml:"host"`
	User     string        `yaml:"user"`
	Password string        `yaml:"password"`
	Timeout  time.Duration `yaml:"timeout"`
}

// ParseConfig parses config from yaml to Config
func ParseConfig() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	cfg := &Config{}

	if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return cfg
}
