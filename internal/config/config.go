package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env" env:"ENV" env-default:"local" env-required="true"`
	StoragePath string `yaml:"storage_path" env-required="true"`
	Postgres    `yaml:"postgres"`
	HTTPServer  `yaml:"http_server"`
}

type Postgres struct {
	Host     string `yaml:"host" env-default:"url-shortener-postgres-db"`
	User     string `yaml:"user" env-required="true"`
	Password string `yaml:"password" env-required="true"`
	DBName   string `yaml:"dbname" env-default:"url_shortener"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
	User        string        `yaml:"user" env-required:"true"`
	Password    string        `yaml:"password" env-required:"true" env:"HTTP_SERVER_PASSWORD"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	//check if file exist
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file %s does not exist", configPath)
	}

	var conf Config

	if err := cleanenv.ReadConfig(configPath, &conf); err != nil {
		log.Fatalf("Cannot read config file: %s", err)
	}

	return &conf
}
