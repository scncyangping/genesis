package shunt

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RunTime struct {
	logg        *zap.SugaredLogger
	mongoConn   *mongo.Client
	mysqlGormDb *gorm.DB
	shuntConfig *Config
}

var (
	runTimeContext = &RunTime{
		shuntConfig: DefaultConfig(),
	}
)

const (
	DefaultDbName = "genesis-shunt" // default db name
	DBTableUser   = "user"          // user table
)

func setLog(l *zap.SugaredLogger) {
	runTimeContext.logg = l
}

func Log() *zap.SugaredLogger {
	return runTimeContext.logg
}

func ShuntConfig() *Config {
	return runTimeContext.shuntConfig
}

func setMongoConn(m *mongo.Client) {
	runTimeContext.mongoConn = m
}

func MongoConn() *mongo.Client {
	return runTimeContext.mongoConn
}

func setGormDb(d *gorm.DB) {
	runTimeContext.mysqlGormDb = d
}

func GormDB() *gorm.DB {
	return runTimeContext.mysqlGormDb
}
