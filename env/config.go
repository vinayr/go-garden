package env

import (
	"github.com/caarlos0/env"
)

// Config represents environment variables
type Config struct {
	// DATABASE
	DbHost     string `env:"POSTGRES_HOST,required"`
	DbName     string `env:"POSTGRES_DB,required"`
	DbUser     string `env:"POSTGRES_USER,required"`
	DbPassword string `env:"POSTGRES_PASSWORD,required"`

	// HTTP Server
	HttpAddr  string `env:"HTTP_ADDR,required"`
	JwtSecret string `env:"JWT_SECRET,required"`
}

// NewConfig returns a new instance of Config
func NewConfig() *Config {
	return &Config{}
}

// Load and parse environment variables
func (c *Config) Load() error {
	if err := env.Parse(c); err != nil {
		return err
	}
	return nil
}
