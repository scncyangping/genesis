package es

import (
	"genesis/pkg/config"

	"github.com/pkg/errors"
)

type Config struct {
	Address  []string
	UserName string
	Password string
}

var _ config.Config = (*Config)(nil)

func (z *Config) Sanitize() {

}

func (z *Config) Validate() error {
	if len(z.Address) < 1 {
		return errors.New("es address is empty")
	}
	if len(z.UserName) < 1 {
		return errors.New("es username is empty")
	}
	if len(z.Password) < 1 {
		return errors.New("es password is empty")
	}
	return nil
}
