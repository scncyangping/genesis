// @Author: YangPing
// @Create: 2023/10/23
// @Description: 服务配置定义

package config

import (
	"genesis/app/common/config"
	"genesis/pkg/plugin/log"
	"genesis/pkg/plugin/mysql"
	"log/slog"
	"os"
)

type ServerConfigI interface {
	GetAddr() string
	GetPort() int
}

type JwtConfigI interface {
	GetSecret() string
	GetExpireTime() int
	GetIssuer() string
	GetAuthKey() string
}

type ServerConfig struct {
	Addr string `yaml:"addr" json:"addr,omitempty"`
	Port int    `yaml:"port" json:"port,omitempty"`
}

type JwtConfig struct {
	Secret     string `yaml:"secret" json:"secret,omitempty"`
	ExpireTime int    `yaml:"expireTime" json:"expireTime,omitempty"`
	Issuer     string `yaml:"issuer" json:"issuer,omitempty"`
	AuthKey    string `yaml:"authKey" json:"authKey,omitempty"`
}

type BaseConfig struct {
	Server *ServerConfig       `yaml:"server" json:"server"` // 服务配置
	Mysql  *config.MysqlConfig `yaml:"mysql" json:"mysql"`   // mysql配置
	Jwt    *JwtConfig          `yaml:"jwt" json:"jwt"`       // jwt配置
	Log    *config.LumberJack  `yaml:"log" json:"log"`       // 日志配置
}

func (c *BaseConfig) Sanitize() {
	return
}

func (c *BaseConfig) Validate() error {
	if err := c.Mysql.Validate(); err != nil {
		return err
	}
	return nil
}

func (c *BaseConfig) Init() error {
	if db, err := mysql.NewMysqlConn(c.Mysql); err != nil {
		return err
	} else {
		runtime.buildGormDb(db)
	}
	var (
		logg = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	)
	// 设置默认日志文件
	if c.Log != nil && c.Log.FileName != "" {
		var (
			logger = log.NewLumberJack(c.Log)
			level  slog.Level
		)
		if c.Log.Level == "DEBUG" {
			level = slog.LevelDebug
		} else {
			level = slog.LevelInfo
		}
		logg = slog.New(slog.NewJSONHandler(logger, &slog.HandlerOptions{
			AddSource:   false,
			Level:       level,
			ReplaceAttr: nil,
		}))
	}
	slog.SetDefault(logg)
	runtime.buildLogger(logg)
	return nil
}

var DefaultBaseConfig = func() *BaseConfig {
	return &BaseConfig{
		Jwt:    &JwtConfig{},
		Server: &ServerConfig{},
		Mysql:  &config.MysqlConfig{},
	}
}

type BaseConfigI interface {
	ServerConfigI
	JwtConfigI
}

func (c *BaseConfig) GetAddr() string {
	return c.Server.Addr
}

func (c *BaseConfig) GetPort() int {
	return c.Server.Port
}

func (c *BaseConfig) GetSecret() string {
	return c.Jwt.Secret
}

func (c *BaseConfig) GetExpireTime() int {
	return c.Jwt.ExpireTime
}

func (c *BaseConfig) GetIssuer() string {
	return c.Jwt.Issuer
}

func (c *BaseConfig) GetAuthKey() string {
	return c.Jwt.AuthKey
}
