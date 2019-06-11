package engine

import (
	"parseExcel/config"
	"parseExcel/reader"
	"fmt"
	"parseExcel/utils"
	"sync"
)

type Engine struct {
	//解析出的sheet数据
	Data SourceData
	//解析出的标题
	Title SourceData
	//数据源Reader
	DataReader reader.Reader
	//解析配置
	configure config.Configure

	blockParser *BlockParser
}

func NewEngine(rd reader.Reader, configure config.Configure) *Engine {
	return &Engine{
		DataReader: rd,
		configure:configure,
		Data: NewSourceData(),
		Title:NewSourceData(),
		blockParser:NewBlockParser(configure),
	}
}

func (e *Engine)Run()  {
	sheets := e.DataReader.GetSheets()
	wg := &sync.WaitGroup{}
	for _, sheetName := range sheets {
		//疑问 []string 不能转 []interface{}
		if !utils.InArray(sheetName, e.configure.Sheets) {
			continue
		}
		e.DataReader.ChangeSheet(sheetName)
		//当前偏移
		rowsData := e.DataReader.Read()
		wg.Add(1)
		e.Schduler(sheetName, rowsData, wg)
	}
	wg.Wait()

}

type DataInterface struct {
	SheetName string
	Data map[string][]interface{}
	Title map[string][]interface{}
}

//schduler
func (e *Engine)Schduler(sheetName string, rowsData [][]interface{}, wg *sync.WaitGroup)  {
	defer wg.Done()
	//当前偏移
	offset := 0
	rowsLen := len(rowsData)

	for offset < rowsLen  {
		fmt.Println(offset)

		var blockData,blockTitle = make(map[string][]interface{}), make(map[string][]interface{})
		//调用解析器
		blockData,blockTitle, offset = e.blockParser.ParserBlock(rowsData, offset)

		//传出数据
		e.Data.Append(sheetName, blockData)
		e.Title.Append(sheetName, blockTitle)
	}
}

func (e *Engine)GetSheetData() SourceData {
	return e.Data
}
func (e *Engine)GetSheetTitle() SourceData {
	return e.Title
}
