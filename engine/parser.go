package engine

import (
	"strconv"
	"parseExcel/config"
	"parseExcel/validate"
)

type BlockParser struct {
	configure config.Configure
}

func NewBlockParser(cfg config.Configure) *BlockParser  {
	return &BlockParser{
		configure:cfg,
	}
}
//返回多块数据
func (parser *BlockParser)ParserBlock(sheetData [][]interface{}, offset int) (
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
		f, title := parser.traitStart(rows, i, sheetData)
		if f {
			currentBlock = title
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
		if parser.traitEnd(rows, i, sheetData) {
			break
		}

		//arrayFilter
		if parser.traitFilter(rows, i, sheetData) {
			continue
		}

		//handle Rows
		blockData[currentBlock] = append(blockData[currentBlock], rows)
	}
	//return blockData.offset.dataTitle
	return blockData, blockTitle, offset + count
}


func (parser *BlockParser)traitStart(rows []interface{}, i int, sheetData[][]interface{}) (bool,string) {
	for funcStr,funcParams := range parser.configure.Block.Start {
		if function, ok := validate.ValidateFuncMap[funcStr]; ok {
			if !function(rows,sheetData,funcParams) {
				return false,strconv.Itoa(i)
			}
		}
	}
	return true,strconv.Itoa(i)
}

func (parser *BlockParser)traitEnd(rows []interface{}, i int, sheetData[][]interface{}) bool {
	for funcStr,funcParams := range parser.configure.Block.Start {
		if function, ok := validate.ValidateFuncMap[funcStr]; ok {
			if !function(rows,sheetData,funcParams) {
				return false
			}
		}
	}
	return true
}

func (parser *BlockParser)traitFilter(rows []interface{}, i int, sheetData[][]interface{}) bool {
	for funcStr,funcParams := range parser.configure.Block.Start {
		if function, ok := validate.ValidateFuncMap[funcStr]; ok {
			if !function(rows,sheetData,funcParams) {
				return false
			}
		}
	}
	return true
}
