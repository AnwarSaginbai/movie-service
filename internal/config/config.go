package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"sync"
)

type Config struct {
	Port    string `yaml:"port"`
	Env     string `yaml:"env"`
	Limiter struct {
		Rps     float64 `yaml:"rps"`
		Burst   int     `yaml:"burst"`
		Enabled bool    `yaml:"enabled"`
	} `yaml:"limiter"`
	Smtp struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Sender   string `yaml:"sender"`
	} `yaml:"mailer"`
	Cors struct {
		TrustedOrigins []string `json:"trusted_origins"`
	} `yaml:"cors" env-default:""`
}

var once sync.Once
var instance *Config

func SetupConfig() *Config {
	once.Do(func() {
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yaml", instance); err != nil {
			log.Fatalln(err)
		}
	})
	return instance
}
