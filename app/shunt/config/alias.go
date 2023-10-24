// @Author: YangPing
// @Create: 2023/10/23
// @Description: 运行时配置

package config

import (
	"gorm.io/gorm"
	"log/slog"
	"sync"
)

var runtime = &RT{config: &BaseConfig{}}

// RT 不允许直接导出
type RT struct {
	gorm   *gorm.DB
	once   sync.Once
	config *BaseConfig
	logger *slog.Logger
}

func Rt() *RT {
	return runtime
}

func (r *RT) Config() BaseConfigI {
	return r.config
}

func (r *RT) GormDB() *gorm.DB {
	return r.gorm
}

func (r *RT) Logger() *slog.Logger {
	return r.logger
}

func (r *RT) BuildConfig(config *BaseConfig) *RT {
	r.once.Do(func() {
		r.config = config
	})
	return r
}

func (r *RT) buildLogger(logger *slog.Logger) *RT {
	r.logger = logger
	return r
}
func (r *RT) buildGormDb(db *gorm.DB) *RT {
	r.gorm = db
	return r
}
