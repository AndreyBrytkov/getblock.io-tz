package models

import "math/big"

type Block struct {
	Number       big.Int
	Transactions []Trasaction
}

type Trasaction struct {
	BlockNumber big.Int `db:"block_num"`
	Idx         big.Int `db:"tx_idx"`
	From        string  `db:"wallet_from"`
	To          string  `db:"wallet_to"`
	Value       big.Int `db:"tx_value"`
	Gas         big.Int `db:"-"`
	GasPrice    big.Int `db:"-"`
	GasUsed     big.Int `db:"-"`
	GasTotal    big.Int `db:"gas_total"`
}
