package models

import "math/big"

type Block struct {
	Number       big.Int
	Transactions []Trasaction
}

type Trasaction struct {
	BlockNumber big.Int
	Idx         big.Int
	From        string
	To          string
	Value       big.Int
	Gas         big.Int
	GasPrice    big.Int
	GasTotal    big.Int
}
