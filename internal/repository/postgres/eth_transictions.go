package postgres

import (
	"errors"
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
	query := `INSERT INTO transactions 
		(block_num, tx_idx, wallet_from, wallet_to, tx_value, gas_total)
	VALUES
		($1, $2, $3, $4, $5, $6);	
	`
	_, err := pg.DB.Exec(query, tx.BlockNumber, tx.Idx, tx.From, tx.To, tx.Value, tx.GasTotal)
	if err != nil {
		return utils.WrapErr(caller, "insert transaction error", err)
	}
	return nil
}

func (pg *Postgres) GetMaxDeltaWallet(n int) (string, big.Int, error) {
	var result big.Int
	return "", result, errors.New("NOT IMPLEMENTED")
}
