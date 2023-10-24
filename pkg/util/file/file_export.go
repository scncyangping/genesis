// @Author: YangPing
// @Create: 2023/10/23
// @Description: Excel文件导出

package file

import (
	"fmt"
	"github.com/samber/lo"
	"github.com/xuri/excelize/v2"
)

type ExcelExport struct {
	header     []*string
	rows       [][]*string
	sheetCount int // row of a sheet
}

func NewExcelExport(header []string, rows [][]*string) *ExcelExport {
	return &ExcelExport{
		header: lo.Map(header, func(item string, _ int) *string {
			return &item
		}),
		rows:       rows,
		sheetCount: 9999,
	}
}

func (e *ExcelExport) WithSheetCount(count int) *ExcelExport {
	e.sheetCount = count
	return e
}

func (e *ExcelExport) Export() (*excelize.File, error) {
	file := excelize.NewFile()
	batchSize := 1000 // 设置每批处理的数据行数
	totalRows := len(e.rows)
	num := 0
	sheetName := "Sheet1" // 使用同一个工作表名
	if sheet, err := file.NewSheet(sheetName); err != nil {
		return nil, err
	} else {
		file.SetActiveSheet(sheet)
	}

	for i2, cell := range e.header {
		cellName := fmt.Sprintf("%c%d", 'A'+int32(i2), 1)
		file.SetCellStr(sheetName, cellName, *cell)
	}

	for offset := 0; offset < totalRows; offset += batchSize {
		end := offset + batchSize
		if end > totalRows {
			end = totalRows
		}
		batchData := e.rows[offset:end]
		if err := e.export(num, sheetName, file, batchData); err != nil {
			return nil, err
		}
		num++
	}
	return file, nil
}

func (e *ExcelExport) export(index int, sheetName string, file *excelize.File, data [][]*string) error {
	c := 'A'
	for rowNum, rows := range data {
		for i2, cell := range rows {
			cellName := fmt.Sprintf("%c%d", c+int32(i2), rowNum+(index*1000)+2)
			file.SetCellStr(sheetName, cellName, *cell)
		}
	}
	return nil
}
