package adapter

import "github.com/AndreyBrytkov/getblock.io-tz/internal/models"


type Repository interface {
	Storage
	GetBlockApi
}

type Storage interface {
	RecordBlock(models.Block) error
	RecordTx(models.Trasaction) error
	GetMaxDeltaWallet(n int) (string, int, error)
}

type GetBlockApi interface {
	GetHeadBlockNum() (uint64, error)
	GetBlockByNum(n uint64) (*models.Block, error)
}
