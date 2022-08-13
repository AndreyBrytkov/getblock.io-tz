package postgres

const (
	SCHEMA_TRANSACTIONS = `--sql
	CREATE TABLE IF NOT EXISTS transactions (
		block_num BIGINT NOT NULL,
		tx_idx INTEGER NOT NULL,
		wallet_from VARCHAR(50) NOT NULL,
		wallet_to VARCHAR(50) NOT NULL, 
		tx_value NUMERIC(32, 0) NOT NULL,
		gas_total NUMERIC(32, 0) NOT NULL,
		CONSTRAINT tx_id PRIMARY KEY (block_num, tx_idx)
	)
	`
)
