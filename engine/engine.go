package engine

import (
	"parseExcel/config"
	"parseExcel/reader"
	"fmt"
	"parseExcel/utils"
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
	fmt.Println(e.configure.Sheets)
	for _, sheetName := range sheets {
		//疑问 []string 不能转 []interface{}
		if !utils.InArray(sheetName, e.configure.Sheets) {
			continue
		}
		e.DataReader.ChangeSheet(sheetName)
		//当前偏移
		rowsData := e.DataReader.Read()
		inputChan := make(chan [][]interface{})
		outChan := make(chan DataInterface)
		go e.Schduler(5, sheetName, inputChan, outChan)
		inputChan <- rowsData
		go func() {
			for dataInterface := range outChan {
				e.Data.Append(dataInterface.SheetName, dataInterface.Data)
				e.Title.Append(dataInterface.SheetName, dataInterface.Title)
			}
		}()

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
				//调用解析器
				blockData,blockTitle, offset = e.blockParser.ParserBlock(rowsData, offset)
				dataInterface.Data = blockData
				dataInterface.Title = blockTitle
				dataInterface.SheetName = sheetName
				//传出数据
				outputChan<-dataInterface
			}
		}(rowsData)
	}
}

func (e *Engine)GetSheetData() SourceData {
	return e.Data
}
func (e *Engine)GetSheetTitle() SourceData {
	return e.Title
}
