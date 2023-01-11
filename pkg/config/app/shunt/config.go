package shunt

import (
	"genesis/pkg/config"
	"genesis/pkg/config/common/log"
	"genesis/pkg/config/common/mongo"
	"genesis/pkg/config/common/mysql"
	"genesis/pkg/config/common/redis"
	logger "genesis/pkg/log"
	_mongo "genesis/pkg/plugin/mongo"
	_mysql "genesis/pkg/plugin/mysql"
	"github.com/pkg/errors"
)

type ServerConfig struct {
	Addr         string `yaml:"addr"`
	ReadTimeout  int    `yaml:"readTimeout"`
	WriteTimeout int    `yaml:"writeTimeout"`
}

var DefaultServerConfig = func() *ServerConfig {
	return &ServerConfig{
		Addr:         ":8080",
		ReadTimeout:  10000,
		WriteTimeout: 10000,
	}
}

type JwtConfig struct {
	Secret     string `yaml:"secret"`
	ExpireTime int    `yaml:"expireTime"`
	Issuer     string `yaml:"issuer"`
	AuthKey    string `yaml:"authKey"`
}

var DefaultJwtConfig = func() *JwtConfig {
	return &JwtConfig{
		Secret:     "secret",
		ExpireTime: 3600,
		Issuer:     "genesis.com",
		AuthKey:    "authorization",
	}
}

type ShuntConfig struct {
	// 服务配置
	Server *ServerConfig
	// 日志配置
	Log *log.LogConfig
	// redis配置
	Redis *redis.RedisConfig
	// mongodb配置
	Mongo *mongo.MongoConfig
	// mysql配置
	Mysql *mysql.MysqlConfig
	// jwt配置
	Jwt *JwtConfig
}

var _ config.Config = (*ShuntConfig)(nil)

func (c *ShuntConfig) Sanitize() {
	c.Log.Sanitize()
	c.Redis.Sanitize()
	c.Mongo.Sanitize()
	c.Mysql.Sanitize()
}

func (c *ShuntConfig) Validate() error {
	err := c.Log.Validate()
	if err != nil {
		return errors.Wrap(err, "log init error")
	}
	err = c.Redis.Validate()
	if err != nil {
		return errors.Wrap(err, "redis init error")
	}
	err = c.Mongo.Validate()
	if err != nil {
		return errors.Wrap(err, "mongo init error")
	}
	err = c.Mysql.Validate()
	if err != nil {
		return errors.Wrap(err, "mysql init error")
	}

	return nil
}

func (c *ShuntConfig) Init() error {
	if conn, err := _mongo.NewMongoConn(c.Mongo); err != nil {
		return err
	} else {
		setMongoConn(conn)
	}

	if newLogger, err := logger.NewLogger(c.Log); err != nil {
		return err
	} else {
		setLog(newLogger)
	}

	conn, err := _mysql.NewMysqlConn(c.Mysql)
	if err != nil {
		return err
	} else {
		setGormDb(conn)
	}
	return nil
}

var DefaultConfig = func() *ShuntConfig {
	return &ShuntConfig{
		Server: DefaultServerConfig(),
		Log:    log.DefaultLogConfig(),
		Redis:  &redis.RedisConfig{},
		Mongo:  &mongo.MongoConfig{},
		Mysql:  &mysql.MysqlConfig{},
		Jwt:    DefaultJwtConfig(),
	}
}
