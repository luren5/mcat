package common

type Transaction struct {
	From     string
	To       string
	Gas      string
	GasPrice string
	Value    string
	Data     string
	Type     uint
}

const (
	TxTypeCommon = iota
	TxTypeContract
)
