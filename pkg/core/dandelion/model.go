package dandelion

type MapBody struct {
	ModuleName string
	Models     []*Model
}

type Model struct {
	ModuleName string
	OldName    string
	DName      string
	XName      string 
	Columns    []*Column `json:"columns,omitempty"`
}

type Column struct {
	DBCloumnDetails `structs:",flatten"`
	ColumnDName     string 
	ColumnXName     string
}

type TableInfo struct {
	DBTableDetails `structs:",flatten"`
	TableDName     string 
	TableXName     string 
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
