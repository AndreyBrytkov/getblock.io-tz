package models

import "math/big"

type Block struct {
	Number       big.Int       `json:"-"`
	NumberHex    string       `json:"number"`
	Transactions []Trasaction `json:"transactions"`
}

type Trasaction struct {
	BlockNumber big.Int `json:"-"`
	Idx         big.Int `json:"-"`
	From        string `json:"from"`
	To          string `json:"to"`
	GasTotal    big.Int `json:"-"`
	Value       big.Int `json:"-"`
	IdxHex      string `json:"transactionIndex"`
	ValueHex    string `json:"value"`
	GasHex      string `json:"gas"`
	GasPriceHex string `json:"gasPrice"`
}
