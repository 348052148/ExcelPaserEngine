package engine

import (
	"parseExcel/config"
	"parseExcel/reader"
	"fmt"
)

type Engine struct {
	//解析出的sheet数据
	Data SourceData
	//解析出的标题
	Title map[string]interface{}
	//数据源Reader
	DataReader reader.Reader
	//解析配置
	configure config.Configure
}

func inArray(need string, needArr []string) bool {
	for _,v := range needArr{
		if need == v{
			return true
		}
	}
	return false
}

func NewEngine(rd reader.Reader, configure config.Configure) *Engine {
	return &Engine{
		DataReader: rd,
		configure:configure,
		Data: NewSourceData(),
		Title:make(map[string]interface{}),
	}
}

func (e *Engine)Run()  {
	sheets := e.DataReader.GetSheets()
	fmt.Println(e.configure.Sheets)
	for _, sheetName := range sheets {
		//疑问 []string 不能转 []interface{}
		if !inArray(sheetName, e.configure.Sheets) {
			continue
		}
		e.DataReader.ChangeSheet(sheetName)
		//当前偏移
		offset := 0
		rowsData := e.DataReader.Read()
		rowsLen := len(rowsData)

		for offset < rowsLen {
			//处理offset数据
			if offset > 0 {
				rowsData = rowsData[offset:]
			}
			var blockData,blockTitle = make(map[string][]interface{}), make(map[string][]interface{})
			blockData,blockTitle, offset = e.ParserBlock(rowsData, offset)
			e.Data.Append(sheetName, blockData)
			e.Title[sheetName] = blockTitle
		}

	}
}

type DataInterface struct {
	SheetName string
	Data map[string][]interface{}
	Title map[string][]interface{}
}


//schduler
func (e *Engine)Schduler(workerCount int, sheetName string, inputChan chan[][]interface{}, outputChan chan DataInterface)  {
	for i:=0; i< workerCount; i++ {
		rowsData := <-inputChan
		go func(rowsData [][]interface{}) {
			//当前偏移
			offset := 0
			rowsLen := len(rowsData)

			dataInterface := DataInterface{
				Data:make(map[string][]interface{}),
				Title:make(map[string][]interface{}),
			}

			for offset < rowsLen {
				//处理offset数据
				if offset > 0 {
					rowsData = rowsData[offset:]
				}
				var blockData,blockTitle = make(map[string][]interface{}), make(map[string][]interface{})
				blockData,blockTitle, offset = e.ParserBlock(rowsData, offset)
				dataInterface.Data = blockData
				dataInterface.Title = blockTitle
				dataInterface.SheetName = sheetName
				//传出数据
				outputChan<-dataInterface
			}
		}(rowsData)
	}
}

//返回多块数据
func (e *Engine)ParserBlock(sheetData [][]interface{}, offset int) (
	map[string][]interface{},
	map[string][]interface{},
	int)  {
	var blockData map[string][]interface{}
	blockData = make(map[string][]interface{})
	//title
	var blockTitle map[string][]interface{}
	blockTitle = make(map[string][]interface{})
	//foreach sheetData
	currentBlock := ""
	count := 0
	for i, rows := range sheetData {
		count++
		//traitTitle
		f, title := traitTitle(rows, i, sheetData)
		if f {
			currentBlock = title
			fmt.Println(currentBlock)
			blockData[currentBlock] = []interface{}{}
			// no handleRows
			blockTitle[currentBlock] = rows
			continue
		}

		//no select
		if currentBlock == "" {
			continue
		}

		//traitEnd
		if traitEnd(rows, i, sheetData) {
			break
		}

		//arrayFilter
		if arrayFilter(rows, i, sheetData) {
			continue
		}

		//handle Rows
		blockData[currentBlock] = append(blockData[currentBlock], rows)
	}
	//return blockData.offset.dataTitle
	return blockData, blockTitle, offset + count
}

func (e *Engine)GetSheetData() SourceData {
	return e.Data
}
