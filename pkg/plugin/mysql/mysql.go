package mysql

import (
	"database/sql"
	myq "genesis/pkg/config/common/mysql"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMysqlConn(cfg *myq.MysqlConfig) (*gorm.DB, error) {
	cs, err2 := cfg.ConnectionString()
	if err2 != nil {
		return nil, err2
	}
	dbCon, err := sql.Open("mysql", cs)

	if err != nil {
		return nil, errors.Wrapf(err, "init mysql error, url:[ %v ]", cs)
	}
	dbCon.SetMaxOpenConns(cfg.MaxOpenConn)
	dbCon.SetMaxIdleConns(cfg.MaxIdleConn)

	err = dbCon.Ping()
	if err != nil {
		return nil, errors.Wrapf(err, "init mysql error, url:[%v]", cs)
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: dbCon,
	}), &gorm.Config{})

	if err != nil {
		return nil, errors.Wrapf(err, "init mysql gormDB error, url:[%v]", cs)
	}

	return gormDB, nil
}
