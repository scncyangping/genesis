// @Author: YangPing
// @Create: 2023/10/21
// @Description: 插件基础配置

package config

import (
	"errors"
	"fmt"
	"genesis/pkg/config"
	"genesis/pkg/types"
	"genesis/pkg/util"
	"strings"
	"time"
)

type EsConfig struct {
	Address  []string
	UserName string
	Password string
	Version  int
}

var _ config.Config = (*EsConfig)(nil)

func (z *EsConfig) Sanitize() {

}

func (z *EsConfig) Validate() error {
	if len(z.Address) < 1 {
		return errors.New("es address is empty")
	}
	if z.Version == 0 || z.Version > 8 {
		return errors.New("es version is error")
	}
	z.UserName = strings.TrimSpace(z.UserName)
	z.Password = strings.TrimSpace(z.Password)
	if z.UserName != "" {
		if d, err := util.AESDecrypt(util.DefaultIv + z.UserName); err != nil {
			return errors.New("user aes decrypt error")
		} else {
			z.UserName = d
		}
	}
	if z.Password != "" {
		if d, err := util.AESDecrypt(util.DefaultIv + z.Password); err != nil {
			return errors.New("password decrypt error")
		} else {
			z.Password = d
		}
	}
	return nil
}

func (z *EsConfig) GetAddress() []string {
	return z.Address
}

func (z *EsConfig) GetUserName() string {
	return z.UserName
}

func (z *EsConfig) GetPassword() string {
	return z.Password
}

func (z *EsConfig) GetVersion() int {
	return z.Version
}

// ZapLogConfig logger
type ZapLogConfig struct {
	LogFileDir    string       `yaml:"logFileDir"` //文件保存地方
	AppName       string       `yaml:"appName"`    //日志文件前缀
	ErrorFileName string       `yaml:"errorFileName"`
	WarnFileName  string       `yaml:"warnFileName"`
	InfoFileName  string       `yaml:"infoFileName"`
	DebugFileName string       `yaml:"debugFileName"`
	MaxSize       int          `yaml:"maxSize"`    //日志文件小大（M）
	MaxAge        int          `yaml:"maxAge"`     //保存的最大天数
	MaxBackups    int          `yaml:"maxBackups"` // 最多存在多少个切片文件
	Mod           types.AppMod `yaml:"Mod"`        //模式 1 正式 2 开发
	Level         string       `yaml:"level"`      //日志等级
	NeedStdout    bool         `yaml:"needStdout"`
	NeedFile      bool         `yaml:"NeedFile"`
}

func (z *ZapLogConfig) Sanitize() {
	if z.LogFileDir == "" {
		z.LogFileDir = fmt.Sprintf("logs%s", types.SP)
	}
	if z.MaxSize == 0 {
		z.MaxSize = 100
	}
	if z.MaxBackups == 0 {
		z.MaxBackups = 60
	}
	if z.MaxAge == 0 {
		z.MaxAge = 30
	}
	if z.Level == "" {
		z.Level = "debug"
	}
	if z.Mod == "" {
		z.Mod = types.DebugMode
	}
}

func (z *ZapLogConfig) Validate() error {
	return nil
}

var DefaultLogConfig = func() ZapLogConfig {
	return ZapLogConfig{
		AppName:       "app",
		ErrorFileName: "err.log",
		WarnFileName:  "warn.log",
		InfoFileName:  "info.log",
		DebugFileName: "debug.log",
		MaxSize:       200,
		MaxBackups:    60,
		MaxAge:        30,
		Level:         "debug", // debug
	}
}

func (z *ZapLogConfig) GetLogFileDir() string {
	return z.LogFileDir
}

func (z *ZapLogConfig) GetAppName() string {
	return z.AppName
}

func (z *ZapLogConfig) GetErrorFileName() string {
	return z.ErrorFileName
}

func (z *ZapLogConfig) GetWarnFileName() string {
	return z.WarnFileName
}

func (z *ZapLogConfig) GetInfoFileName() string {
	return z.InfoFileName
}

func (z *ZapLogConfig) GetDebugFileName() string {
	return z.DebugFileName
}

func (z *ZapLogConfig) GetMaxSize() int {
	return z.MaxSize
}

func (z *ZapLogConfig) GetMaxAge() int {
	return z.MaxAge
}

func (z *ZapLogConfig) GetMaxBackups() int {
	return z.MaxBackups
}

func (z *ZapLogConfig) GetMod() types.AppMod {
	return z.Mod
}

func (z *ZapLogConfig) GetLevel() string {
	return z.Level
}

func (z *ZapLogConfig) GetNeedStdout() bool {
	return z.NeedStdout
}

func (z *ZapLogConfig) GetNeedFile() bool {
	return z.NeedFile
}

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

func (p *MysqlConfig) GetMaxOpenConn() int {
	return p.MaxOpenConn
}

func (p *MysqlConfig) GetMaxIdleConn() int {
	return p.MaxIdleConn
}

func (p *MysqlConfig) GetConnectionString() (string, error) {
	return p.ConnectionString()
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

	if d, err := util.AESDecrypt(util.DefaultIv + p.User); err != nil {
		return errors.New("user aes decrypt error")
	} else {
		p.User = d
	}

	if len(p.Password) < 1 {
		return errors.New("password should not be empty")
	}

	if d, err := util.AESDecrypt(util.DefaultIv + p.Password); err != nil {
		return errors.New("password decrypt error")
	} else {
		p.Password = d
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
