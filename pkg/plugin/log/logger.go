// @Author: YangPing
// @Create: 2023/10/23
// @Description: zap logger 插件配置

package log

import (
	"fmt"
	"genesis/pkg/types"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"time"
)

const DefaultLogName = "log.log"

type Config interface {
	GetLogFileDir() string //文件保存地方
	GetAppName() string    //日志文件前缀
	GetErrorFileName() string
	GetWarnFileName() string
	GetInfoFileName() string
	GetDebugFileName() string
	GetMaxSize() int      //日志文件小大（M）
	GetMaxAge() int       //保存的最大天数
	GetMaxBackups() int   // 最多存在多少个切片文件
	GetMod() types.AppMod //模式 1 正式 2 开发
	GetLevel() string     //日志等级
	GetNeedStdout() bool  // 控制台输出
	GetNeedFile() bool    // 文件方式记录
}

func NewLogger(z Config) (*zap.SugaredLogger, error) {
	logger, err := newLogger(z)
	if err != nil {
		return nil, err
	}
	return logger.Sugar(), nil
}

func newLogger(z Config) (*zap.Logger, error) {
	var zc zap.Config
	var encoder zapcore.Encoder
	if z.GetMod() != types.ReleaseMode {
		zc = zap.NewDevelopmentConfig()
		zc.EncoderConfig.EncodeTime = timeEncoder

		encoderConfig := zap.NewDevelopmentEncoderConfig()
		encoderConfig.EncodeTime = timeEncoder
		//encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		zc = zap.NewProductionConfig()
		zc.EncoderConfig.EncodeTime = timeUnixNano

		encoder = zapcore.NewJSONEncoder(zc.EncoderConfig)
	}

	zc.Level.SetLevel(zapcore.Level(types.LogLevelM[z.GetLevel()]))

	path, err := filepath.Abs(z.GetLogFileDir())
	if err != nil {
		return nil, err
	}
	fn := func(fileName string) zapcore.WriteSyncer {
		if fileName == "" {
			fileName = DefaultLogName
		}
		fp := fmt.Sprintf("%s%s%s-%s", path, types.SP, z.GetAppName(), fileName)

		// 同时输出控制台和文件
		syncWriters := make([]zapcore.WriteSyncer, 0)
		if z.GetNeedStdout() {
			syncWriters = append(syncWriters, zapcore.AddSync(os.Stdout))
		}

		if z.GetNeedFile() {
			syncWriters = append(syncWriters, zapcore.AddSync(&lumberjack.Logger{
				Filename:   fp,
				MaxSize:    z.GetMaxSize(),
				MaxBackups: z.GetMaxBackups(),
				MaxAge:     z.GetMaxAge(),
				Compress:   true,
				LocalTime:  true,
			}))
		}
		return zapcore.NewMultiWriteSyncer(syncWriters...)
	}

	errPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.ErrorLevel && zapcore.ErrorLevel-zapcore.Level(types.LogLevelM[z.GetLevel()]) > -1
	})
	warnPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.WarnLevel && zapcore.WarnLevel-zapcore.Level(types.LogLevelM[z.GetLevel()]) > -1
	})
	infoPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.InfoLevel && zapcore.InfoLevel-zapcore.Level(types.LogLevelM[z.GetLevel()]) > -1
	})
	debugPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.DebugLevel && zapcore.DebugLevel-zapcore.Level(types.LogLevelM[z.GetLevel()]) > -1
	})

	cores := []zapcore.Core{
		zapcore.NewCore(encoder, fn(z.GetErrorFileName()), errPriority),
		zapcore.NewCore(encoder, fn(z.GetWarnFileName()), warnPriority),
		zapcore.NewCore(encoder, fn(z.GetInfoFileName()), infoPriority),
		zapcore.NewCore(encoder, fn(z.GetDebugFileName()), debugPriority),
	}

	return zc.Build(zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return zapcore.NewTee(cores...)
	}))
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}
func timeUnixNano(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendInt64(t.UnixNano() / 1e6)
}
