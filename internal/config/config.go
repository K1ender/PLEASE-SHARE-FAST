package config

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTP HTTP
}

type HTTP struct {
	Port int    `env:"HTTP_PORT" env-default:"8080"`
	Addr string `env:"HTTP_ADDR" env-default:"0.0.0.0"`
}

func MustInit() Config {
	_, err := os.Stat(".env")
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}

	var cfg Config

	if os.IsNotExist(err) {
		err := cleanenv.ReadEnv(&cfg)
		if err != nil {
			panic(err)
		}
	} else {
		err := cleanenv.ReadConfig(".env", &cfg)
		if err != nil {
			panic(err)
		}
	}

	return cfg
}
