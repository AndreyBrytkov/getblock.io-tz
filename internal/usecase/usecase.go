package usecase

import (
  "math/big"

  "github.com/AndreyBrytkov/getblock.io-tz/internal/adapter"
  "github.com/AndreyBrytkov/getblock.io-tz/internal/models"
  "github.com/AndreyBrytkov/getblock.io-tz/pkg/utils"
)

const caller = "usecase"

// Main package obgect
type Usecase struct {
  logger  *utils.MyLogger
  config  *models.AppConfig
  api     adapter.GetBlockApi
  storage adapter.Storage
}

// Getter
func GetUsecase(logger *utils.MyLogger, config *models.AppConfig, api adapter.GetBlockApi, storage adapter.Storage) adapter.Usecase {
  return &Usecase{
    logger:  logger,
    config:  config,
    api:     api,
    storage: storage,
  }
}

func (u *Usecase) GetMaxBalanceDeltaWallet() (string, *big.Int, error) {
  // Get chain head block number
  lastBlockNum, err := u.api.GetHeadBlockNum()
  if err != nil {
    return "", nil, utils.WrapErr(caller, "get chain head block number error", err)
  }

  // Substract 100 blocks to get number where start calculations from
  startBlock := big.NewInt(0)
  startBlock.Sub(&lastBlockNum, big.NewInt(int64(u.config.BlockAmount)))

  // Get all transactions needed for calculations
  transactions, err := u.storage.GetTransactionsByBlocksRange(*startBlock, lastBlockNum)
  if err != nil {
    return "", nil, utils.WrapErr(caller, "get transactions error", err)
  }

  // Calculate balance changes for all wallets
  balance := make(map[string]*big.Int)
  for _, tx := range transactions {
    // Change sender balance
    balanceSender, ok := balance[tx.From]
    if !ok {
      balanceSender = big.NewInt(0)
    }
    balanceSender.Sub(balanceSender, &tx.Value)
    balanceSender.Sub(balanceSender, &tx.GasTotal)
    balance[tx.From] = balanceSender // not nesessary?

    // Change reciever balance
    balanceReciever, ok := balance[tx.From]
    if !ok {
		balanceReciever = big.NewInt(0)
    }
    balanceReciever.Add(balanceReciever, &tx.Value)
    balance[tx.From] = balanceReciever // not nesessary?
  }

  // Find wallet with max delta
  maxDelta := big.NewInt(0)
  walletWithMaxDelta := ""

  for wallet, balanceChange := range balance {
    if balanceChange.CmpAbs(maxDelta) > 0 {
      maxDelta = balanceChange
      walletWithMaxDelta = wallet
    }
  }

  return walletWithMaxDelta, maxDelta, nil
}