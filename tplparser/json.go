package tplparser

import (
	"encoding/json"
	"os"
	"parseExcel/config"
)

type JsonParser struct {
	Decoder *json.Decoder
	Data config.Configure
}

func NewJsonParser(file *os.File) *JsonParser {
	return &JsonParser{
		Decoder:json.NewDecoder(file),
	}
}

func (parser *JsonParser)Parse() interface{} {
	err := parser.Decoder.Decode(&parser.Data)
	if err != nil {
		panic(err)
	}
	return parser.Data
}
