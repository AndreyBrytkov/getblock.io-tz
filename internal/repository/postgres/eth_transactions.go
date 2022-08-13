package postgres

import (
	"fmt"
	"math/big"

	"github.com/AndreyBrytkov/getblock.io-tz/internal/models"
	"github.com/AndreyBrytkov/getblock.io-tz/pkg/utils"
)

func (pg *Postgres) RecordBlock(block models.Block) error {
	for _, tx := range block.Transactions {
		err := pg.RecordTx(tx)
		if err != nil {
			return utils.WrapErr(caller, "Record transaction error", err)
		}
	}
	return nil
}

func (pg *Postgres) RecordTx(tx models.Trasaction) error {
	query := `--sql
	INSERT INTO transactions 
		(block_num, tx_idx, wallet_from, wallet_to, tx_value, gas_total)
	VALUES
		($1, $2, $3, $4, $5, $6)
	ON CONFLICT DO NOTHING;	
	`
	_, err := pg.DB.Exec(query, tx.BlockNumber.String(), tx.Idx.String(), tx.From, tx.To, tx.Value.String(), tx.GasTotal.String())
	if err != nil {
		return utils.WrapErr(caller, "insert transaction error", err)
	}
	return nil
}

func (pg *Postgres) GetTransactionsByBlocksRange(from, to big.Int) ([]models.Trasaction, error) {
	query := `--sql
	SELECT block_num, tx_idx, wallet_from, wallet_to, tx_value, gas_total
	FROM transactions
	WHERE block_num  >= $1 AND block_num  <= $2;  
	`

	trasactions := []TrasactionSqlx{}
	err := pg.DB.Select(&trasactions, query, from.String(), to.String())
	if err != nil {
		errMsg := fmt.Sprintf("get transactions by blocks range (from '%s' to '%s') error", from.String(), to.String())
		return nil, utils.WrapErr(caller, errMsg, err)
	}

	result, err := pg.convertTxs(trasactions)
	if err != nil {
		return nil, utils.WrapErr(caller, "convert transactions error", err)
	}

	return result, nil
}

func (pg *Postgres) convertTxs(txsSqlx []TrasactionSqlx) ([]models.Trasaction, error) {
	result := []models.Trasaction{}

	for _, tx := range txsSqlx {
		bNum := new(big.Int)
		_, ok := bNum.SetString(tx.BlockNumber, 10)
		if !ok {
			return nil, fmt.Errorf("convert block number string '%s' to big.Int error", tx.BlockNumber)
		}

		idx := new(big.Int)
		_, ok = idx.SetString(tx.Idx, 10)
		if !ok {
			return nil, fmt.Errorf("convert transaction index string '%s' to big.Int error", tx.Idx)
		}

		val := new(big.Int)
		_, ok = val.SetString(tx.Value, 10)
		if !ok {
			return nil, fmt.Errorf("convert transaction value string '%s' to big.Int error", tx.Value)
		}

		gas := new(big.Int)
		_, ok = gas.SetString(tx.GasTotal, 10)
		if !ok {
			return nil, fmt.Errorf("convert transaction gas string '%s' to big.Int error", tx.GasTotal)
		}

		newTx := models.Trasaction{
			BlockNumber: *bNum,
			Idx:         *idx,
			From:        tx.From,
			To:          tx.To,
			Value:       *val,
			GasTotal:    *gas,
		}
		result = append(result, newTx)
	}

	return result, nil
}

type TrasactionSqlx struct {
	BlockNumber string `db:"block_num"`
	Idx         string `db:"tx_idx"`
	From        string `db:"wallet_from"`
	To          string `db:"wallet_to"`
	Value       string `db:"tx_value"`
	Gas         string `db:"-"`
	GasPrice    string `db:"-"`
	GasUsed     string `db:"-"`
	GasTotal    string `db:"gas_total"`
}
