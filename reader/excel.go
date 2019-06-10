package reader

import "github.com/360EntSecGroup-Skylar/excelize"

type ExcelReader struct {
	file *excelize.File
	selectSheet string
}

func NewExcelReader(filePath string) *ExcelReader {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return &ExcelReader{}
	}
	return &ExcelReader{
		file: f,
	}
}

func (reader *ExcelReader)GetSheets() map[int]string  {
	return reader.file.GetSheetMap()
}

func (reader *ExcelReader)ChangeSheet(sheetName string)  {
	reader.selectSheet = sheetName
}

func (reader *ExcelReader)Read() [][]interface{} {
	rows, err := reader.file.GetRows(reader.selectSheet)
	if err != nil {
		panic(err)
	}
	// 设置单元格的值
	var sheetData [][]interface{}
	for _, row := range rows {
		var rowData []interface{}
		for _, colCell := range row {
			rowData = append(rowData, colCell)
		}
		sheetData = append(sheetData, rowData)
	}
	return sheetData
}