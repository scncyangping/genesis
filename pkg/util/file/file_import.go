// @Author: YangPing
// @Create: 2023/10/23
// @Description: Excel文件导入

package file

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"
	"io"
	"os"
	"strings"
)

const Split = "@~@"

type fetchByExcel struct {
	// the Excel file reader
	excelReader io.Reader
	// the Excel file path
	excelPath string
	// data row start index in the Excel
	rowNumIndex int
	// target column index in the Excel
	//columnIndex []int
	// max import num
	importNum int

	// row num check
	//importCheckOp []ban.ImportCheckOp
}

func NewFetchByExcel() *fetchByExcel {
	return &fetchByExcel{
		rowNumIndex: 2,
		importNum:   10000, // 默认10000
	}
}

//
//func (e *fetchByExcel) BuildImportCheckOp(options ...ban.ImportCheckOp) *fetchByExcel {
//	e.importCheckOp = append(e.importCheckOp, options...)
//	return e
//}
//
//func (e *fetchByExcel) BuildColumnIndex(num []int) *fetchByExcel {
//	sort.Ints(num)
//	e.columnIndex = num
//	return e
//}

func (e *fetchByExcel) BuildRowNumIndex(num int) *fetchByExcel {
	e.rowNumIndex = num
	return e
}

func (e *fetchByExcel) BuildEndImportNum(num int) *fetchByExcel {
	e.importNum = num
	return e
}

func (e *fetchByExcel) BuildExcelPath(excelPath string) *fetchByExcel {
	e.excelPath = excelPath
	return e
}

func (e *fetchByExcel) BuildExcelReader(excelReader io.Reader) *fetchByExcel {
	e.excelReader = excelReader
	return e
}

func (e *fetchByExcel) Sanitize() error {
	if e.rowNumIndex == 0 {
		e.rowNumIndex = 2
	}

	if e.importNum == 0 {
		e.importNum = 10000
	}

	if e.importNum > 10000 {
		return errors.New("导入数据不能超过10000条!")
	}
	if e.excelPath == "" && e.excelReader == nil {
		return errors.New("an excelPath or excelReader must exist")
	}

	return nil
}

func (e *fetchByExcel) FetchKeyWord() ([]*[]string, error) {
	var (
		f   *excelize.File
		err error
	)

	if err := e.Sanitize(); err != nil {
		return nil, err
	}

	if e.excelPath != "" {
		f, err = excelize.OpenFile(e.excelPath)
	} else {
		f, err = excelize.OpenReader(e.excelReader)
	}

	if err != nil {
		return nil, errors.Wrap(err, "import excel error")

	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			fmt.Printf("import excel close file error: %s", err.Error())
		}
	}()

	sheetList := f.GetSheetList()

	needItems := make([]*[]string, 0)

	for _, sheetName := range sheetList {
		if arr, err := e.do(sheetName, f); err != nil {
			return needItems, err
		} else {
			needItems = append(needItems, arr...)
		}
	}
	return needItems, err
}

func (e *fetchByExcel) do(sheetName string, f *excelize.File) ([]*[]string, error) {
	var needItems []*[]string

	rows, err := f.Rows(sheetName)

	defer rows.Close()

	if err != nil {
		return needItems, err
	}
	cur := 0
	for rows.Next() {
		cur++
		if cur < e.rowNumIndex {
			continue
		}
		if cur >= e.rowNumIndex+e.importNum {
			break
		}
		row, err := rows.Columns()
		if err != nil {
			break
		}
		//if (len(row) < len(e.columnIndex)) || (e.columnIndex[len(e.columnIndex)-1] > len(row)) {
		//	return nil, errors.New("导入数据列数和需要传入列数不匹配!")
		//}

		item := make([]string, 0)
		for _, value := range row {
			//if len(e.columnIndex) > 0 {
			//	if lo.Contains(e.columnIndex, index) {
			//		item = append(item, value)
			//	}
			//} else {
			//	item = append(item, value)
			//}
			item = append(item, value)
		}
		needItems = append(needItems, &item)
	}
	return needItems, nil
}

type fetchByTxt struct {
	// the Excel file reader
	txtReader io.Reader
	// the Excel file path
	txtPath   string
	importNum int
}

func NewFetchByTxt() *fetchByTxt {
	return &fetchByTxt{}
}

func (e *fetchByTxt) BuildImportNum(num int) *fetchByTxt {
	e.importNum = num
	return e
}

func (e *fetchByTxt) BuildTextPath(txtPath string) *fetchByTxt {
	e.txtPath = txtPath
	return e
}

func (e *fetchByTxt) BuildTextReader(textReader io.Reader) *fetchByTxt {
	e.txtReader = textReader
	return e
}

func (e *fetchByTxt) Sanitize() error {
	if e.importNum == 0 {
		e.importNum = 100000
	}
	if e.txtPath == "" && e.txtReader == nil {
		return errors.New("an txtPath or txtReader must exist")
	}

	return nil
}

func (e *fetchByTxt) FetchKeyWord() ([]string, error) {
	var (
		f   *os.File
		err error
	)

	if err := e.Sanitize(); err != nil {
		return nil, err
	}

	if e.txtPath != "" {
		f, err = os.OpenFile(e.txtPath, os.O_RDONLY, 0666)
		e.txtReader = f
	}

	if err != nil {
		return nil, errors.Wrap(err, "import txt error")
	}

	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			fmt.Printf("import txt close file error: %s ", err.Error())
		}
	}()

	reader := bufio.NewReader(e.txtReader)
	needTitles := make([]string, 0)

	for i := 0; ; i++ {
		if i >= e.importNum {
			break
		}
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, errors.Wrap(err, "import txt readline error")
		}
		needTitles = append(needTitles, strings.TrimSpace(string(line)))
	}

	return needTitles, nil
}
