package engine

func traitTitle(rows []interface{}, i int, sheetData[][]interface{}) (bool,string) {
	if rows[0] == "Strategy" {
		return true, "SS"
	}
	return false,"1"
}

func traitEnd(rows []interface{}, i int, sheetData[][]interface{}) bool {
	if rows[0] == "HONG KONG" {
		return true
	}
	return false
}

func arrayFilter(rows []interface{}, i int, sheetData[][]interface{}) bool {

	return false
}
