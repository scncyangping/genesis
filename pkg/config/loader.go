package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"os"
)

func Load(file string, cfg Config) error {
	return LoadWithOption(file, cfg, false, true, true)
}

func LoadWithOption(file string, cfg Config, strict bool, includeEnv bool, validate bool) error {
	if file == "" {
		fmt.Println("Skipping reading config from file")
	} else if err := loadFromFile(file, cfg, strict); err != nil {
		return err
	}

	if includeEnv {
		if err := envconfig.Process("", cfg); err != nil {
			return err
		}
	}
	if validate {
		if err := cfg.Validate(); err != nil {
			return errors.Wrapf(err, "Invalid configuration")
		}
	}
	if c, ok := cfg.(InitConfig); ok {
		err := c.(InitConfig).Init()
		if err != nil {
			return err
		}
	}
	return nil
}

func loadFromFile(file string, cfg Config, strict bool) error {
	if _, err := os.Stat(file); err != nil {
		return errors.Errorf("Failed to access configuration file %q", file)
	}
	contents, err := os.ReadFile(file)
	if err != nil {
		return errors.Wrapf(err, "Failed to read configuration from file %q", file)
	}
	if strict {
		err = yaml.UnmarshalStrict(contents, cfg)
	} else {
		err = yaml.Unmarshal(contents, cfg)
	}
	return errors.Wrapf(err, "Failed to parse configuration from file %q", file)
}
