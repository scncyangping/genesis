package redis

import (
	"errors"
	"genesis/pkg/config"
	"time"
)

type RedisConfig struct {
	PoolSize     int           `yaml:"poolSize"`
	Addr         []string      `yaml:"addr"`
	Pwd          string        `yaml:"pwd"`
	DialTimeout  time.Duration `yaml:"DialTimeout"`
	ReadTimeout  time.Duration `yaml:"ReadTimeout"`
	WriteTimeout time.Duration `yaml:"WriteTimeout"`
}

var _ config.Config = (*RedisConfig)(nil)

func (p *RedisConfig) Sanitize() {

}

func (p *RedisConfig) Validate() error {
	if len(p.Addr) < 1 {
		return errors.New("addr should not be empty")
	}
	return nil
}
