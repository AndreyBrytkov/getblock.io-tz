package models

type Block struct {
	Number       uint64
	NumberHex    string       `json:"number"`
	Transactions []Trasaction `json:"transactions"`
}

type Trasaction struct {
	BlockNumber uint64
	From        string `json:"from"`
	To          string `json:"to"`
	ValueHex    string `json:"value"`
	Value       uint64
}
