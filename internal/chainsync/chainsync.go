package chainsync

import (
	"math/big"
	"sync"
	"time"

	"github.com/AndreyBrytkov/getblock.io-tz/internal/adapter"
	"github.com/AndreyBrytkov/getblock.io-tz/internal/models"
	"github.com/AndreyBrytkov/getblock.io-tz/pkg/utils"
)

const caller = "chainsync"

// Main package obgect
type ChainSynchronizer struct {
	// client          *jsonrpc.Client
	logger          *utils.MyLogger
	config          *models.AppConfig
	wg              *sync.WaitGroup
	api             adapter.GetBlockApi
	storage         adapter.Storage
	lastLoadedBlock big.Int
	lastSync        int64
}

// Getter
func GetChainSynchronizer(l *utils.MyLogger, config *models.AppConfig, wg *sync.WaitGroup, api adapter.GetBlockApi, storage adapter.Storage) *ChainSynchronizer {
	return &ChainSynchronizer{
		logger:  l,
		config:  config,
		wg:      wg,
		api:     api,
		storage: storage,
	}
}

func (cs *ChainSynchronizer) Init() error {
	cs.logger.Info(caller, "initializtion started...")

	// Get lastest block number
	lastest, err := cs.api.GetHeadBlockNum()
	if err != nil {
		return utils.WrapErr(caller, "get lastest block num error", err)
	}

	blockAmount := big.NewInt(int64(cs.config.BlockAmount))
	cs.lastLoadedBlock.Sub(&lastest, blockAmount)
	cs.lastSync = time.Now().Unix() - 10

	cs.logger.Info(caller, "successfully initiated...")
	return nil
}

func (cs *ChainSynchronizer) Run() {
	cs.logger.Info(caller, "sync started...")
	defer cs.wg.Done()

	for {
		if time.Now().Unix() < cs.lastSync + int64(cs.config.Cycle) {
			continue
		}
		// Get lastest block number
		lastest, err := cs.api.GetHeadBlockNum()
		if err != nil {
			err = utils.WrapErr(caller, "get lastest block num error", err)
			cs.logger.Fatal(err)
		}
		cs.lastSync = time.Now().Unix()

		blocksToLoad := []big.Int{}
		if lastest.Cmp(&cs.lastLoadedBlock) > 0 {
			blocksToLoad = utils.GetBlockNumsToLoad(cs.lastLoadedBlock, lastest)
			cs.lastLoadedBlock = lastest
		}
		cs.logger.Debug(caller, "sync %d block", len(blocksToLoad))
		localWg := &sync.WaitGroup{}

		for _, blockNum := range blocksToLoad {
			localWg.Add(1)
			go func(wg *sync.WaitGroup, n big.Int) {
				defer wg.Done()
				// Get block
				block, err := cs.api.GetBlockByNum(n)
				if err != nil {
					err = utils.WrapErr(caller, "get block error", err)
					cs.logger.Fatal(err)
					return
				}

				// Save block in storage
				err = cs.storage.RecordBlock(*block)
				if err != nil {
					err = utils.WrapErr(caller, "save block error", err)
					cs.logger.Fatal(err)
					return
				}
				cs.logger.Info(caller, "block '%s' saved in DB", n.String())
			}(localWg, blockNum)
		}
		localWg.Wait()
	}
}
