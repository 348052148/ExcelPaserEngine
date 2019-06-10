package reader

type Reader interface {
	Read() [][]interface{}
	GetSheets() map[int]string
	ChangeSheet(sheetName string)
}

