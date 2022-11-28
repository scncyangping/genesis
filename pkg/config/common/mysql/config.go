package mysql

import (
	"errors"
	"fmt"
	"genesis/pkg/config"
	"strings"
)

const (
	// Mysql 最大空闲连接数
	_defaultMysqlMaxIdleConn = 1000
	// Mysql 最大连接数
	_defaultMysqlMaxOpenConn = 2000
	// 默认还是连接
	_defConnectInfo = "charset=utf8&parseTime=true"
)

type MysqlConfig struct {
	// 连接用户名
	User string
	// 连接密码
	Password string
	// 连接地址
	Host string
	// 连接地址数组 集群连接时使用
	Hosts []string
	// 数据库名称
	DbName string `yaml:"dbName"`
	// 连接参数,配置编码、是否启用ssl等 charset=utf8
	ConnectInfo string `yaml:"connectInfo"`
	// 用于设置最大打开的连接数，默认值为0表示不限制
	MaxOpenConn int `yaml:"maxOpenConn"`
	// 用于设置闲置的连接数
	MaxIdleConn int `yaml:"maxIdleConn"`
	// 连接名称 用于多个连接时区分
	ConnName string
}

var _ config.Config = (*MysqlConfig)(nil)

func (p *MysqlConfig) Sanitize() {
	if p.MaxOpenConn == 0 {
		p.MaxOpenConn = _defaultMysqlMaxOpenConn
	}
	if p.MaxIdleConn == 0 {
		p.MaxIdleConn = _defaultMysqlMaxIdleConn
	}
	if p.ConnectInfo == "" {
		p.ConnectInfo = _defConnectInfo
	}
}

func (p *MysqlConfig) Validate() error {
	if len(p.Host) < 1 && len(p.Hosts) < 1 {
		return errors.New("host should not be empty")
	}
	if len(p.User) < 1 {
		return errors.New("user should not be empty")
	}
	if len(p.Password) < 1 {
		return errors.New("password should not be empty")
	}
	if len(p.DbName) < 1 {
		return errors.New("DbName should not be empty")
	}
	return nil
}

func (p *MysqlConfig) ConnectionString() (string, error) {
	escape := func(value string) string { return strings.ReplaceAll(strings.ReplaceAll(value, `\`, `\\`), `'`, `\'`) }

	return fmt.Sprintf("%s:%s@tcp(%s)/%s?%s",
		escape(p.User),
		escape(p.Password),
		escape(p.Host),
		escape(p.DbName),
		escape(p.ConnectInfo)), nil
}
