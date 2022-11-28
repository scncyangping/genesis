package logger

import (
	"fmt"
	"genesis/pkg/config/common/log"
	"genesis/pkg/types"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"path/filepath"
	"time"
)

func NewLogger(z *log.LogConfig) (*zap.SugaredLogger, error) {
	logger, err := newLogger(z)
	if err != nil {
		return nil, err
	}

	return logger.Sugar(), nil
}

func newLogger(z *log.LogConfig) (*zap.Logger, error) {
	var zc zap.Config
	var encoder zapcore.Encoder
	if z.Mod != types.ReleaseMode {
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

	zc.Level.SetLevel(zapcore.Level(types.LogLevelM[z.Level]))

	path, err := filepath.Abs(z.LogFileDir)
	if err != nil {
		return nil, err
	}
	fn := func(fileName string) zapcore.WriteSyncer {
		if fileName == "" {
			fileName = "log.log"
		}
		fp := fmt.Sprintf("%s%s%s-%s", path, types.SP, z.AppName, fileName)

		return zapcore.AddSync(&lumberjack.Logger{
			Filename:   fp,
			MaxSize:    int(z.MaxSize),
			MaxBackups: z.MaxBackups,
			MaxAge:     int(z.MaxAge),
			Compress:   true,
			LocalTime:  true,
		})
	}

	errPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.ErrorLevel && zapcore.ErrorLevel-zapcore.Level(types.LogLevelM[z.Level]) > -1
	})
	warnPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.WarnLevel && zapcore.WarnLevel-zapcore.Level(types.LogLevelM[z.Level]) > -1
	})
	infoPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.InfoLevel && zapcore.InfoLevel-zapcore.Level(types.LogLevelM[z.Level]) > -1
	})
	debugPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.DebugLevel && zapcore.DebugLevel-zapcore.Level(types.LogLevelM[z.Level]) > -1
	})

	cores := []zapcore.Core{
		zapcore.NewCore(encoder, fn(z.ErrorFileName), errPriority),
		zapcore.NewCore(encoder, fn(z.WarnFileName), warnPriority),
		zapcore.NewCore(encoder, fn(z.InfoFileName), infoPriority),
		zapcore.NewCore(encoder, fn(z.DebugFileName), debugPriority),
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
