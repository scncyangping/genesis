package dandelion

type MapBody struct {
	ModuleName string
	Models     []*Model
}

type Model struct {
	ModuleName string
	OldName    string // 数据库表名 下划线
	DName      string // 大驼峰
	XName      string // 小驼峰
	Columns    []*Column
}

type Column struct {
	DBCloumnDetails
	ColumnDName string // 大驼峰
	ColumnXName string // 小驼峰
}

type TableInfo struct {
	DBTableDetails
	TableDName string // 表名大驼峰
	TableXName string // 表名小驼峰
}

type DBCloumnDetails struct {
	ColumnName    string `gorm:"column:columnName"`
	DataType      string `gorm:"column:dataType"`
	ColumnComment string `gorm:"column:columnComment"`
	ColumnKey     string `gorm:"column:columnKey"`
	Extra         string `gorm:"column:extra"`
}

type DBTableDetails struct {
	TableName    string `gorm:"column:tableName"`
	Engine       string `gorm:"column:engine"`
	TableComment string `gorm:"column:columnName"`
	CreateTime   string `gorm:"column:dataType"`
}
