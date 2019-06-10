package config

type Configure struct {
	Engine string
	Mode string
	Sheets []string
	Block Block
	Filter Filter
}

type Block struct {
	Start []string
	Ends []string
	Title string
}

type Filter struct {
	Rules []string
}

type RowFormat struct {
	kvs map[string]string
}

