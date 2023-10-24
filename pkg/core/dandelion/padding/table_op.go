package padding

import (
	"fmt"
	"genesis/pkg/config/app/shunt"
	"genesis/pkg/core/dandelion"
	"genesis/pkg/util"

	"github.com/fatih/structs"
	"gorm.io/gorm"
)

// TablePadding is used to create padding within a table
type TablePadding struct {
	moduleName string
	tables     []string
	tableOp    *TableOp
}

func NewTablePadding(moduleName string, tables []string, tableOp *TableOp) *TablePadding {
	return &TablePadding{moduleName: moduleName, tables: tables, tableOp: tableOp}
}

func (ta *TablePadding) Add() map[string]any {
	m := make(map[string]any)
	if ta.tables == nil || len(ta.tables) < 1 {
		return m
	}

	data := dandelion.MapBody{
		ModuleName: ta.moduleName,
		Models:     make([]*dandelion.Model, 0),
	}

	for _, v := range ta.tables {
		model := dandelion.Model{
			ModuleName: ta.moduleName,
			OldName:    v,
			DName:      util.Camelize(v, true),
			XName:      util.Camelize(v, false),
			Columns:    ta.tableOp.Columns(v),
		}
		data.Models = append(data.Models, &model)
	}

	return structs.Map(data)
}

// 默认mysql数据库字段类型与golang内置类型对应关系
var mapMixture = map[string]string{
	"bigint":    "string",
	"tinyint":   "int",
	"varchar":   "string",
	"int":       "int",
	"decimal":   "float64",
	"datetime":  "time.Time",
	"timestamp": "int64",
	"time":      "int64",
	"text":      "string",
	"json":      "string",
}

type TableOp struct {
	db           string
	gm           *gorm.DB
	fieldMapping map[string]string
}

func NewTableOp(db string) *TableOp {
	return &TableOp{
		db:           db,
		fieldMapping: mapMixture,
	}
}

func (t *TableOp) BuildGormDB(g *gorm.DB) *TableOp {
	t.gm = g
	return t
}

func (t *TableOp) BuildMatchMap(m map[string]string) *TableOp {
	t.fieldMapping = m
	return t
}

func (t *TableOp) Tables() []*dandelion.TableInfo {
	qSql := `SELECT
		 table_name as tableName,
		 engine as engine,
		 table_comment as tableComment,
		 create_time as createTime
		 FROM
		 information_schema.tables
		 WHERE table_schema = (?)`
	tableInfos := make([]*dandelion.TableInfo, 0)

	t.gm.Raw(qSql, t.db).Scan(&tableInfos)

	for _, v := range tableInfos {
		v.TableDName = util.Camelize(v.TableName, true)
		v.TableXName = util.Camelize(v.TableName, false)
		fmt.Printf("%v \n", v)
	}
	return tableInfos

}
func (t *TableOp) Columns(tableName string) []*dandelion.Column {
	qSql := `SELECT
		column_name as columnName,
		data_type as dataType,
		column_comment as columnComment,
		column_key as columnKey,
		extra
		FROM information_schema.columns 
		WHERE table_name = ? and table_schema = (?) order by ordinal_position`

	columns := make([]*dandelion.Column, 0)

	shunt.GormDB().Raw(qSql, tableName, t.db).Scan(&columns)

	for _, v := range columns {
		v.ColumnDName = util.Camelize(v.ColumnName, true)
		v.ColumnXName = util.Camelize(v.ColumnName, false)
		v.DataType = t.fieldMapping[v.DataType]
		fmt.Printf("%v \n", v)
	}
	return columns
}
