package mongo

import (
	"fmt"
	"genesis/pkg/config"
	"github.com/pkg/errors"
	"strings"
)

const (
	_defaultUri = "mongodb://%s"
)

type MongoConfig struct {
	Host            string
	User            string
	DbName          string `yaml:"dbName"`
	Password        string
	PoolSize        uint64 `yaml:"poolSize"`
	MaxConnIdleTime uint64 `yaml:"maxConnIdleTime"`
}

var _ config.Config = (*MongoConfig)(nil)

func (p *MongoConfig) Sanitize() {

}

func (p *MongoConfig) Validate() error {
	if p.Host == "" {
		return errors.New("mongo host should not be empty")
	}

	if len(p.User) != 0 {
		if len(p.User) < 1 {
			return errors.New("mongo user should not be empty")
		}
		if len(p.Password) < 1 {
			return errors.New("mongo password should not be empty")
		}
	}

	return nil
}

func (p *MongoConfig) ConnectionString() (string, error) {
	escape := func(value string) string { return strings.ReplaceAll(strings.ReplaceAll(value, `\`, `\\`), `'`, `\'`) }

	return fmt.Sprintf(_defaultUri, escape(p.Host)), nil
}
