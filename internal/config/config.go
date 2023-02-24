package config

import (
	"io/ioutil"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gopkg.in/yaml.v2"
)

type Config struct {
	ServerPort    int    `yaml:"server_port" env:"SERVER_PORT"`
	DSN           string `yaml:"dsn" env:"DSN,secret"`
	JWKS          string `yaml:"jwks" env:"JWKS,secret"`
	KetoRead      string `yaml:"keto_read" env:"KetoRead,secret"`
	KetoWrite     string `yaml:"keto_write" env:"KetoWrite,secret"`
	API           string `yaml:"api" env:"API,secret"`
	RedisEndpoint string `yaml:"redis_endpoint" env:"REDIS_ENDPOINT"`
	RedisPassword string `yaml:"redis_password" env:"REDIS_PASSWORD"`
	Mongo         string `yaml:"mongo" env:"Mongo,secret"`
}

// Validate the configuration values.
func (c Config) Validate() error {

	return validation.ValidateStruct(&c,
		validation.Field(&c.DSN, validation.Required),
		validation.Field(&c.JWKS, validation.Required),
		validation.Field(&c.KetoRead, validation.Required),
		validation.Field(&c.KetoRead, validation.Required),
		validation.Field(&c.API, validation.Required),
		validation.Field(&c.RedisEndpoint, validation.Required),
		validation.Field(&c.RedisPassword, validation.Required),
		validation.Field(&c.Mongo, validation.Required),
	)
}

// Load the configuration from a file.
func Load(file string) (*Config, error) {

	c := Config{}
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	if err = yaml.Unmarshal(bytes, &c); err != nil {
		return nil, err
	}
	if err = c.Validate(); err != nil {
		return nil, err
	}
	return &c, err
}
