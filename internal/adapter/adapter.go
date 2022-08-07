package adapter

import (
	"math/big"

	"github.com/AndreyBrytkov/getblock.io-tz/internal/models"
)

type Storage interface {
	RecordBlock(models.Block) error
	RecordTx(models.Trasaction) error
	GetMaxDeltaWallet(n int) (string, big.Int, error)
}

type GetBlockApi interface {
	GetHeadBlockNum() (big.Int, error)
	GetBlockByNum(n big.Int) (*models.Block, error)
}
