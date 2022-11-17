package config

import (
	"io/ioutil"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gopkg.in/yaml.v2"
)

const (
	defaultServerPort         = 8080
	defaultJWTExpirationHours = 72
)

type Config struct {
	ServerPort    int    `yaml:"server_port" env:"SERVER_PORT"`
	DSN           string `yaml:"dsn" env:"DSN,secret"`
	JWTSigningKey string `yaml:"jwt_signing_key" env:"JWT_SIGNING_KEY,secret"`
	JWTExpiration int    `yaml:"jwt_expiration" env:"JWT_EXPIRATION"`
}

func (c Config) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.DSN, validation.Required),
	)
}

func Load(file string) (*Config, error) {
	c := Config{
		ServerPort:    defaultServerPort,
		JWTExpiration: defaultJWTExpirationHours,
	}
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
