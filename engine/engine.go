package engine

import (
	"parseExcel/config"
	"parseExcel/reader"
	"fmt"
	"parseExcel/utils"
	"parseExcel/schduler"
)

type Engine struct {
	//解析出的sheet数据
	Data *SourceData
	//解析出的标题
	Title *SourceData
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
	//数据通道
	dataMessageChan := make(chan DataMessage, 10)
	//数据生成
	//调度器
	sd1 := schduler.NewSchduler()
	for _, sheetName := range sheets {
		//疑问 []string 不能转 []interface{}
		if !utils.InArray(sheetName, e.configure.Sheets) {
			continue
		}
		e.DataReader.ChangeSheet(sheetName)
		sd1.AddTask(func() error {
			//当前偏移
			rowsData := e.DataReader.Read()
			dataMessageChan <- DataMessage{
				SheetName:sheetName,
				Data:rowsData,
			}
			fmt.Println("sd")
			return nil
		})
	}
	sd1.Start()
	close(dataMessageChan)
	//调度器2 由于依赖之前的数据，不能放到同一个调度器里
	sd2 := schduler.NewSchduler()
	for i:=0; i < 3; i++ {
		sd2.AddTask(func() error {
			for {
				select {
				case dataMessage, ok := <-dataMessageChan:
					if !ok {
						return nil
					}
					e.Parse(dataMessage.SheetName, dataMessage.Data)
				}
			}
			return nil
		})
	}
	sd2.Start()

}

//schduler
func (e *Engine)Parse(sheetName string, rowsData [][]interface{})  {
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

func (e *Engine)GetSheetData() *SourceData {
	return e.Data
}
func (e *Engine)GetSheetTitle() *SourceData {
	return e.Title
}
