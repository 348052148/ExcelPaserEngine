package validate

import (
	"strconv"
)

var ValidateFuncMap map[string]ValidateFunc

func init()  {
	ValidateFuncMap = make(map[string]ValidateFunc)
	ValidateFuncMap["EqualColVal"] = EqualColVal
	ValidateFuncMap["RowNextColVal"] = RowNextColVal
}
type ValidateFunc func([]interface{}, int, [][]interface{},interface{}) bool

func EqualColVal(rows []interface{}, index int, sheetRows [][]interface{}, parames interface{}) bool  {
	p:=parames.([]interface{})
	i,_ := strconv.Atoi(p[0].(string))
	if rows[i] == p[1].(string) {
		return true
	}
	return false
}

func RowNextColVal(rows []interface{}, index int,sheetRows [][]interface{}, parames interface{}) bool {
	p:=parames.([]interface{})
	i,_ := strconv.Atoi(p[0].(string))
	if index+1 >= len(sheetRows) {
		return false
	}
	if sheetRows[index+1][i] == p[1].(string) {
		return true
	}
	return false
}
