package shunt

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RunTime struct {
	logger      *zap.SugaredLogger
	mongoConn   *mongo.Client
	mysqlGormDb *gorm.DB
	shuntConfig *ShuntConfig
}

func NewRunTime() *RunTime {
	return &RunTime{}
}

func (r *RunTime) buildLogger(logger *zap.SugaredLogger) *RunTime {
	r.logger = logger
	return r
}

func (r *RunTime) buildMongoConn(mc *mongo.Client) *RunTime {
	r.mongoConn = mc
	return r
}

func (r *RunTime) buildGormDb(gd *gorm.DB) *RunTime {
	r.mysqlGormDb = gd
	return r
}

func (r *RunTime) buildConfig(config *ShuntConfig) *RunTime {
	r.shuntConfig = config
	return r
}

var (
	runTimeContext = NewRunTime().buildConfig(DefaultConfig())
)

func MongoConn() *mongo.Client {
	return runTimeContext.mongoConn
}

func Config() *ShuntConfig {
	return runTimeContext.shuntConfig
}

func Log() *zap.SugaredLogger {
	return runTimeContext.logger
}

func GormDB() *gorm.DB {
	return runTimeContext.mysqlGormDb
}
