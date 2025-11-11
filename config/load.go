package config

import (
	_ "embed"
	"os"
	"regexp"

	"github.com/go-playground/validator/v10"
	logger "github.com/rtfmkiesel/kisslog"
	"go.yaml.in/yaml/v3"
)

var log = logger.New("config")

func LoadFromFile(configPath string) (*Config, error) {
	log.Debug("Reading config from '%s'", configPath)

	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		return nil, log.NewError(err)
	}

	log.Debug("Unmarshaling config into *Config")
	c := &Config{}
	if err := yaml.Unmarshal(configBytes, &c); err != nil {
		return nil, log.NewError(err)
	}

	log.Debug("Validating config")
	if err := c.Validate(); err != nil {
		return nil, log.NewError(err)
	}

	log.Info("Config loaded & validated")
	return c, nil
}

func (c *Config) Validate() error {
	validate := validator.New()

	if err := validate.RegisterValidation("regex", validateRegex); err != nil {
		return log.NewError("could not register regex validation: %s", err)
	}

	return validate.Struct(c)
}

// Checks if a string field compiles to a regex pattern
func validateRegex(fl validator.FieldLevel) bool {
	regexStr := fl.Field().String()
	_, err := regexp.Compile(regexStr)
	return err == nil
}
