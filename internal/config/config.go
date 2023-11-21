package config

import (
	"io/ioutil"
	"reflect"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Endpoint struct {
		Management string `yaml:"management" env:"Management"`
		Check_REST string `yaml:"check_rest" env:"Check_REST"`
		Check_GRPC string `yaml:"check_grpc" env:"Check_GRPC"`
	} `yaml:"endpoint"`
	Auth struct {
		JWKS string `yaml:"jwks" env:"JWKS"`
	} `yaml:"auth"`
	Database struct {
		URL      string `yaml:"url" env:"URL,secret"`
		Name     string `yaml:"name" env:"Name,secret"`
		User     string `yaml:"user" env:"User,secret"`
		Password string `yaml:"password" env:"Password,secret"`
	} `yaml:"database"`
	RootOrganization struct {
		Name            string `yaml:"name" env:"Name"`
		AdminIdentifier string `yaml:"admin_identifier" env:"AdminIdentifier"`
		AdminName       string `yaml:"admin_name" env:"AdminName"`
		AdminRoleName   string `yaml:"admin_role_name" env:"AdminRoleName"`
	} `yaml:"root_organization"`
	SystemResources struct {
		Organizations []string `yaml:"organizations"`
		Users         []string `yaml:"users"`
		Roles         []string `yaml:"roles"`
		Groups        []string `yaml:"groups"`
		Resources     []string `yaml:"resources"`
		Polices       []string `yaml:"policies"`
	} `yaml:"system_resources"`
	APIEndpoints []APIEndpoint `yaml:"endpoints"`
}

type APIEndpoint struct {
	Path     string         `yaml:"path"`
	Methods  []MethodDetail `yaml:"methods"`
	Resource string         `yaml:"resource"`
}

type MethodDetail struct {
	Method              string   `yaml:"method"`
	RequiredPermissions []string `yaml:"required_permissions"`
}

func Nested(target interface{}, fieldRules ...*validation.FieldRules) *validation.FieldRules {
	return validation.Field(target, validation.By(func(value interface{}) error {
		valueV := reflect.Indirect(reflect.ValueOf(value))
		if valueV.CanAddr() {
			addr := valueV.Addr().Interface()
			return validation.ValidateStruct(addr, fieldRules...)
		}
		return validation.ValidateStruct(target, fieldRules...)
	}))
}

// Validate the configuration values.
func (c Config) Validate() error {

	return validation.ValidateStruct(&c,
		Nested(&c.Endpoint,
			validation.Field(&c.Endpoint.Management, validation.Required),
			validation.Field(&c.Endpoint.Check_REST, validation.Required),
		),
		Nested(&c.Database,
			validation.Field(&c.Database.URL, validation.Required),
			validation.Field(&c.Database.Name, validation.Required),
			validation.Field(&c.Database.User, validation.Required),
			validation.Field(&c.Database.Password, validation.Required),
		),
		Nested(&c.RootOrganization,
			validation.Field(&c.RootOrganization.Name, validation.Required),
		),
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
