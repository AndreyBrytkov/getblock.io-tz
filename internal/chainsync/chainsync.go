package chainsync

import (
	"math/big"
	"sync"
	"time"

	"github.com/AndreyBrytkov/getblock.io-tz/internal/adapter"
	"github.com/AndreyBrytkov/getblock.io-tz/internal/models"
	"github.com/AndreyBrytkov/getblock.io-tz/pkg/utils"
	"github.com/ubiq/go-ubiq/common/hexutil"
)

const caller = "chainsync"

// Main package obgect
type ChainSynchronizer struct {
	logger          *utils.MyLogger
	config          *models.AppConfig
	wg              *sync.WaitGroup
	api             adapter.GetBlockApi
	storage         adapter.Storage
	lastLoadedBlock big.Int
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

	cs.logger.Info(caller, "successfully initiated...")
	return nil
}

func (cs *ChainSynchronizer) Run() {
	cs.logger.Info(caller, "sync started...")
	defer cs.wg.Done()

	for {
		// Get lastest block number
		lastest, err := cs.api.GetHeadBlockNum()
		if err != nil {
			err = utils.WrapErr(caller, "get lastest block num error", err)
			cs.logger.Fatal(err)
		}

		blocksToLoad := []big.Int{}
		if lastest.Cmp(&cs.lastLoadedBlock) > 0 {
			cs.logger.Debug(caller, "calculating blocks between %s and %s", hexutil.EncodeBig(&cs.lastLoadedBlock), hexutil.EncodeBig(&lastest))
			blocksToLoad = getBlockNumsToLoad(lastest, cs.lastLoadedBlock)
			cs.lastLoadedBlock = lastest
		}
		cs.logger.Debug(caller, "%d blocks to load", len(blocksToLoad))

		for _, blockNum := range blocksToLoad {
			// Get block
			block, err := cs.api.GetBlockByNum(blockNum)
			if err != nil {
				err = utils.WrapErr(caller, "get block error", err)
				cs.logger.Fatal(err)
				goto exit
			}

			// Save block in storage
			err = cs.storage.RecordBlock(*block)
			if err != nil {
				err = utils.WrapErr(caller, "save block error", err)
				cs.logger.Fatal(err)
				goto exit
			}
		}

		// Wait
		time.Sleep(time.Duration(cs.config.Cycle) * time.Second)
	}
exit:
	cs.logger.Info(caller, "sync stopped!")
}

func getBlockNumsToLoad(lastest, lastLoaded big.Int) []big.Int {
	numList := []big.Int{}

	// newBlockNum := lastLoaded + 1
	newBlockNum := big.NewInt(0)
	newBlockNum.Add(&lastLoaded, big.NewInt(1))

	// while newBlockNum <= lastest
	for newBlockNum.Cmp(&lastest) <= 0 {
		// fmt.Println("[DEBUG] --:--:-- [CHAINSYNC]: new block to load %d", newBlockNum)
		numList = append(numList, *newBlockNum)
		// newBlockNum++
		newBlockNum.Add(newBlockNum, big.NewInt(1))
	}

	return numList
}
