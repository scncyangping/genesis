// @Author: YangPing
// @Create: 2023/11/30
// @Description: 日志分割插件

package log

import "gopkg.in/natefinch/lumberjack.v2"

type LumberJackI interface {
	GetFileName() string // 日志文件位置
	GetMaxSize() int     // 文件大小MB
	GetMaxBackups() int  // 保留旧文件最大数量
	GetMaxAge() int      // 保留旧文件最大天数
	GetCompress() bool   // 是否压缩
}

func NewLumberJack(config LumberJackI) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   config.GetFileName(),
		MaxSize:    config.GetMaxSize(),
		MaxAge:     config.GetMaxAge(),
		MaxBackups: config.GetMaxBackups(),
		Compress:   config.GetCompress(),
	}
}
