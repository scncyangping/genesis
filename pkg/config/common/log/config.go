package log

import (
	"fmt"
	"genesis/pkg/config"
	"genesis/pkg/types"
)

var DefaultLogConfig = func() *LogConfig {
	return &LogConfig{
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

// LogConfig logger
type LogConfig struct {
	LogFileDir    string       `yaml:"logFileDir"` //文件保存地方
	AppName       string       `yaml:"appName"`    //日志文件前缀
	ErrorFileName string       `yaml:"errorFileName"`
	WarnFileName  string       `yaml:"warnFileName"`
	InfoFileName  string       `yaml:"infoFileName"`
	DebugFileName string       `yaml:"debugFileName"`
	MaxSize       uint8        `yaml:"maxSize"`    //日志文件小大（M）
	MaxAge        uint8        `yaml:"maxAge"`     //保存的最大天数
	MaxBackups    int          `yaml:"maxBackups"` // 最多存在多少个切片文件
	Mod           types.AppMod `yaml:"Mod"`        //模式 1 正式 2 开发
	Level         string       `yaml:"level"`      //日志等级
}

var _ config.Config = (*LogConfig)(nil)

func (z *LogConfig) Sanitize() {
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

func (z LogConfig) Validate() error {
	return nil
}
