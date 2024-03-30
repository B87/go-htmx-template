package internal

import "github.com/kelseyhightower/envconfig"

type Config struct {
	ServerPort int `envconfig:"SERVER_PORT" default:"8080"`

	PGHost     string `envconfig:"PG_HOST" default:"localhost"`
	PGPort     int    `envconfig:"PG_PORT" default:"5432"`
	PGUser     string `envconfig:"PG_USER" default:"postgres"`
	PGPassword string `envconfig:"PG_PASSWORD" default:"postgres"`

	CDNBucket string `envconfig:"CDN_BUCKET" default:""`

	// Logging
	LogLevel string `envconfig:"LOG_LEVEL" default:"info" enum:"debug,info"`
}

func MustNewConfig() *Config {
	cfg := &Config{}
	err := envconfig.Process("", cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}
