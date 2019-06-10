package engine

import (
	"encoding/json"
	"bytes"
)

type SourceData struct {
	data map[string] map[string][]interface{}
}

func NewSourceData() SourceData  {
	return SourceData{
		data:make(map[string] map[string][]interface{}),
	}
}

func (s SourceData)ToJson() string {
	byt,_ := json.Marshal(s.data)
	return bytes.NewBuffer(byt).String()
}

func (s SourceData)Set(sheet string, data map[string][]interface{})  {
	s.data[sheet] = data
}

func (s SourceData)Append(sheet string, data map[string][]interface{})  {
	if _,ok := s.data[sheet]; !ok {
		s.Set(sheet, data)
	}else {
		for k,v := range data {
			s.data[sheet][k] = v
		}
	}
}
