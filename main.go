package main

import (
	"flag"
	"os"
	"parseExcel/tplparser"
	"fmt"
	config2 "parseExcel/config"
	"parseExcel/engine"
	"parseExcel/reader"
)

func main() {
	filePath := flag.String("file", "/Users/zhouhui/go/src/parseExcel/excel/庆瑞所有投资组合20190402.xlsx", "--文件路径")
	tpl := flag.String("tpl", "/Users/zhouhui/go/src/parseExcel/tpls/default.json", "--解析模版")
	//目标：
	// 1. 解析引擎
	// 2. tpl 解析方式 - yaml json prototype
	// 3. 解析驱动，- go-excelize 或者其他解析引擎
	// 4. validate - func
	// 5. sourceData
	// 6. row format
	//实现：
	file, err := os.Open(*tpl)
	if err != nil {
		panic(err)
	}
	parser := tplparser.NewJsonParser(file)
	config := parser.Parse()
	fmt.Println(config.(config2.Configure).Sheets)
	eg := engine.NewEngine(reader.NewExcelReader(*filePath), config.(config2.Configure))
	eg.Run()
	fmt.Println(eg.GetSheetData().ToJson())
	//fmt.Println(eg.Title)
}
