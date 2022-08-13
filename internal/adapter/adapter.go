package adapter

import (
	"math/big"

	"github.com/AndreyBrytkov/getblock.io-tz/internal/models"
)

type Rest interface {
	Run()
}

type Usecase interface {
	GetMaxBalanceDeltaWallet() (string, *big.Int, error)
}

type Storage interface {
	RecordBlock(models.Block) error
	RecordTx(models.Trasaction) error
	GetTransactionsByBlocksRange(from, to big.Int) ([]models.Trasaction, error)
}

type GetBlockApi interface {
	GetHeadBlockNum() (big.Int, error)
	GetBlockByNum(n big.Int) (*models.Block, error)
}
