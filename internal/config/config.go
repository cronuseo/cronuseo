package config

import (
	"io/ioutil"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gopkg.in/yaml.v2"
)

type Config struct {
	JWKS          string `yaml:"jwks" env:"JWKS,secret"`
	Mgt_API       string `yaml:"mgt_api" env:"Mgt_API,secret"`
	Check_API     string `yaml:"check_api" env:"Check_API,secret"`
	Mongo         string `yaml:"mongo" env:"Mongo,secret"`
	MongoUser     string `yaml:"mongouser" env:"MongoUser,secret"`
	MongoPassword string `yaml:"mongopassword" env:"MongoPassword,secret"`
	MongoDBName   string `yaml:"mongodbname" env:"MongoDBName,secret"`
	DefaultOrg    string `yaml:"default_org" env:"DEFAULT_ORG"`
	RBACPolicy    string `yaml:"rbac_policy" env:"RBACPolicy"`
}

// Validate the configuration values.
func (c Config) Validate() error {

	return validation.ValidateStruct(&c,
		validation.Field(&c.JWKS, validation.Required),
		validation.Field(&c.Mgt_API, validation.Required),
		validation.Field(&c.Check_API, validation.Required),
		validation.Field(&c.MongoUser, validation.Required),
		validation.Field(&c.MongoPassword, validation.Required),
		validation.Field(&c.DefaultOrg, validation.Required),
		validation.Field(&c.MongoDBName, validation.Required),
		validation.Field(&c.RBACPolicy, validation.Required),
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
