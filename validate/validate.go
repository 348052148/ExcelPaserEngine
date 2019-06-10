package validate

import (
	"strconv"
)

var ValidateFuncMap map[string]ValidateFunc

func init()  {
	ValidateFuncMap = make(map[string]ValidateFunc)
	ValidateFuncMap["EqualColVal"] = EqualColVal
}
type ValidateFunc func([]interface{}, [][]interface{},interface{}) bool

func EqualColVal(rows []interface{}, sheetRows [][]interface{}, parames interface{}) bool  {
	p:=parames.([]interface{})
	i,_ := strconv.Atoi(p[0].(string))
	if rows[i] == p[1].(string) {
		return true
	}
	return false
}
