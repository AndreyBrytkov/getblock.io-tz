package postgres

const (
	SCHEMA_TRANSACTIONS = `CREATE TABLE IF NOT EXISTS transactions (
		block_num BIGINT NOT NULL,
		tx_idx BIGINT NOT NULL,
		wallet_from VARCHAR(50) NOT NULL,
		wallet_to VARCHAR(50) NOT NULL, 
		tx_value BIGINT NOT NULL,
		gas_total BIGINT NOT NULL,
		CONSTRAINT tx_id PRIMARY KEY (block_num, tx_idx)
	)
	`
)
