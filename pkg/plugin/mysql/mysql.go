package mysql

import (
	"database/sql"
	"fmt"
	myq "genesis/pkg/config/common/mysql"
	"genesis/pkg/core/shunt/repository/entity"
	"reflect"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"
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

// BuildGormQuery
// en: struct、map
// rz: remove zero filed
func BuildGormQuery(bod any, rz bool, tx *gorm.DB) (*gorm.DB, error) {
	t := reflect.TypeOf(bod)
	switch t.Kind() {
	case reflect.Struct:
		return buildStruct(bod, rz, tx)
	case reflect.Map:
		return buildMap(bod, rz, tx)
	}
	return tx, nil
}

var TimeQueryGormBuild = []string{"create_at", "update_at"}

func RegisterTimeQueryBuild(s []string) {
	TimeQueryGormBuild = append(TimeQueryGormBuild, s...)
}

func buildMap(en any, rz bool, tx *gorm.DB) (*gorm.DB, error) {
	for key, v := range en.(map[string]any) {
		t := reflect.TypeOf(v)
		vv := reflect.ValueOf(v)
		iz := vv.IsZero()
		// 不是零值或者是零值但是不移除零值
		if !(!iz || (iz && !rz)) {
			continue
		}
		switch t.Kind() {
		case reflect.Slice:
			if lo.Contains(TimeQueryGormBuild, key) && len(v.([]any)) == 2 {
				if v.([]any)[0] != "" {
					tx = tx.Where(fmt.Sprintf("%s >= ?", key), v.([]any)[0])
				}
				if v.([]any)[1] != "" {
					tx = tx.Where(fmt.Sprintf("%s <= ?", key), v.([]any)[1])
				}
			} else {
				tx = tx.Where(fmt.Sprintf("%s IN ?", key), v)
			}
		case reflect.String:
			if strings.Contains(key, "like") {
				if strings.Contains(v.(string), "%") {
					tx = tx.Where(fmt.Sprintf("%s '%s'", key, v))
				} else {
					tx = tx.Where(fmt.Sprintf("%s '%s'", key, v.(string)+"%"))
				}
			} else {
				tx = tx.Where(fmt.Sprintf("%s = ?", key), v)
			}
		default:
			tx = tx.Where(fmt.Sprintf("%s = ?", key), v)
		}
	}
	return tx, nil
}

func buildStruct(en any, rz bool, tx *gorm.DB) (*gorm.DB, error) {
	t := reflect.TypeOf(en)
	v := reflect.ValueOf(en)
	for k := 0; k < t.NumField(); k++ {
		cv := v.Field(k).Interface()
		iz := v.Field(k).IsZero()
		// 不是零值或者是零值但是不移除零值
		if !(!iz || (iz && !rz)) {
			continue
		}

		if _, ok := cv.(entity.En); ok {
			buildStruct(cv, true, tx)
			continue
		}
		// use tag
		var cn string
		for _, v := range strings.Split(t.Field(k).Tag.Get("gorm"), ";") {
			cl := strings.Split(v, ":")
			if len(cl) != 2 {
				continue
			}
			if cl[0] == "column" {
				cn = cl[1]
				break
			}
		}
		if cn == "" {
			continue
		}
		// check
		//pm[cn] = cv
		ft := v.Field(k).Kind()
		switch ft {
		case reflect.Array:
			tx = tx.Where(fmt.Sprintf("%s IN ?", cn), cv)
		default:
			tx = tx.Where(fmt.Sprintf("%s = ?", cn), cv)
		}
	}
	return tx, nil
}
