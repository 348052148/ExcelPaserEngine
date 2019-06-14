package engine

import (
	"encoding/json"
	"bytes"
	"sync"
)

type SourceData struct {
	data map[string] map[string][]interface{}
	lock *sync.RWMutex
}

func NewSourceData() *SourceData  {
	return &SourceData{
		data:make(map[string] map[string][]interface{}),
		lock:&sync.RWMutex{},
	}
}

func (s *SourceData)ToJson() string {
	s.lock.RLock()
	s.lock.RUnlock()
	byt,_ := json.Marshal(s.data)
	return bytes.NewBuffer(byt).String()
}

func (s *SourceData)Set(sheet string, data map[string][]interface{})  {
	s.data[sheet] = data
}

func (s *SourceData)Append(sheet string, data map[string][]interface{})  {
	s.lock.Lock()
	defer s.lock.Unlock()
	if _,ok := s.data[sheet]; !ok {
		s.Set(sheet, data)
	}else {
		for k,v := range data {
			s.data[sheet][k] = v
		}
	}
}


type DataMessage struct {
	SheetName string
	Data [][]interface{}
}
