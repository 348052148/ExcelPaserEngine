package engine

import (
	"strconv"
	"parseExcel/config"
	"parseExcel/validate"
	"fmt"
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
	for i:=offset; i < len(sheetData); i++ {
		rows := sheetData[i]
	//for i, rows := range sheetData[offset:] {
		count++
		//块数据开始
		f, title := parser.traitStart(rows, i, sheetData)
		if f {
			currentBlock = title
			blockData[currentBlock] = []interface{}{}
			// no handleRows
			blockTitle[currentBlock] = rows
			continue
		}
		//未找到块开始时，跳过数据行
		if currentBlock == "" {
			continue
		}
		//块数据结尾
		if parser.traitEnd(rows, i, sheetData) {
			break
		}
		//过滤满足条件的数据行
		if parser.traitFilter(rows, i, sheetData) {
			continue
		}
		//处理rows数据
		blockData[currentBlock] = append(blockData[currentBlock], rows)
	}
	//return blockData.offset.dataTitle
	return blockData, blockTitle, offset + count
}
//验证开始
func (parser *BlockParser)traitStart(rows []interface{}, i int, sheetData[][]interface{}) (bool,string) {
	//默认为自动生成的标题
	title := "$"+strconv.Itoa(i)
	//获取标题
	if parser.configure.Block.Title != "" {
		tidx,_:= strconv.Atoi(parser.configure.Block.Title)
		title = rows[tidx].(string)
	}

	if parser.RulesValidate(parser.configure.Block.Start, rows, i, sheetData) {
		return true,title
	}

	return false,title
}

func (parser *BlockParser)RulesValidate(rules map[string]interface{}, rows []interface{}, i int, sheetData[][]interface{}) bool  {
	if rules == nil {
		return false
	}
	for funcStr,funcParams := range rules {
		if function, ok := validate.ValidateFuncMap[funcStr]; ok {
			if !function(rows, i, sheetData,funcParams) {
				return false
			}
		}
	}
	return true
}

//验证结尾
func (parser *BlockParser)traitEnd(rows []interface{}, i int, sheetData[][]interface{}) bool {
	if parser.RulesValidate(parser.configure.Block.Ends, rows, i, sheetData) {
		return true
	}
	return false
}
//过滤验证
func (parser *BlockParser)traitFilter(rows []interface{}, i int, sheetData[][]interface{}) bool {
	if parser.RulesValidate(parser.configure.Filter.Rules, rows, i, sheetData) {
		return true
	}
	return false
}
