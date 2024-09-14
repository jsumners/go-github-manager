package config

import (
	"errors"
	"fmt"
	"github.com/creasty/defaults"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type Config struct {
	*viper.Viper `mapstructure:"-"`
	AuthToken    string `mapstructure:"auth_token"`
	LogLevel     string `mapstructure:"log_level" default:"info"`
}

func New() *Config {
	cfg := &Config{
		Viper: viper.New(),
	}
	_ = defaults.Set(cfg)
	return cfg
}

func (c *Config) ReadConfig(configFile string) error {
	c.SetConfigName(".ghm")
	c.SetConfigType("yaml")
	c.AddConfigPath(".")
	c.AddConfigPath("$HOME")
	c.SetEnvPrefix("GHM")
	c.AutomaticEnv()

	envFile := c.GetString("config_file")
	if configFile != "" {
		// This path means that `--config-file <file>` flag has been supplied.
		// Therefore, we want to prefer it over the environment.
		c.SetConfigFile(configFile)
	} else if envFile != "" {
		// Fallback to the file specified by `GHM_CONFIG_FILE`.
		c.SetConfigFile(envFile)
	}

	var configFileNotFoundError viper.ConfigFileNotFoundError
	err := c.ReadInConfig()
	if err != nil && errors.As(err, &configFileNotFoundError) == false {
		return fmt.Errorf("unable to read configuration file: %w", err)
	}

	err = c.Unmarshal(c)
	if err != nil {
		return fmt.Errorf("unable to unmarshal configuration: %w", err)
	}

	// Frustratingly, AutomaticEnv() does not read in environment variables
	// when ReadInConfig() is invoked. So we have to prime it ourselves.
	c.AuthToken = c.GetString("auth_token")

	return nil
}

func (c *Config) GenerateCurrentYaml() (string, error) {
	encodedData := map[string]any{}
	err := mapstructure.Decode(c, &encodedData)
	if err != nil {
		return "", fmt.Errorf("unable to encode configuration: %w", err)
	}

	data, err := yaml.Marshal(encodedData)
	if err != nil {
		return "", fmt.Errorf("unable to marshal to yaml: %w", err)
	}

	return string(data), nil
}

func (c *Config) GenerateDefaultYaml() (string, error) {
	defaultConfig := Config{}
	err := defaults.Set(&defaultConfig)
	if err != nil {
		return "", fmt.Errorf("unable to generate default config: %w", err)
	}

	// We decode to a generic interface in order to rename the struct fields
	// according to the `mapstructure` tag.
	var encoded map[string]any
	err = mapstructure.Decode(defaultConfig, &encoded)
	if err != nil {
		return "", fmt.Errorf("unable to decode configuration: %w", err)
	}

	yamlEncoded, err := yaml.Marshal(encoded)
	if err != nil {
		return "", fmt.Errorf("unable to encode configuration to yaml: %w", err)
	}

	return string(yamlEncoded), nil
}
