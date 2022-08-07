package main

import (
	"sync"

	"github.com/AndreyBrytkov/getblock.io-tz/internal/chainsync"
	"github.com/AndreyBrytkov/getblock.io-tz/internal/repository"
	"github.com/AndreyBrytkov/getblock.io-tz/pkg/utils"
)

const caller = "app"

func main() {
	wg := new(sync.WaitGroup)
	defer wg.Wait()

	logger := utils.NewLogger()

	config, err := utils.LoadConfig()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info(caller, "config loaded...")
	logger.DebugOn = config.AppConfig.Debug

	repo := repository.GetRepository(logger, config)

	cs := chainsync.GetChainSynchronizer(logger, &config.AppConfig, wg, repo.Api, repo.Storage)

	err = cs.Init()
	if err != nil {
		err = utils.WrapErr(caller, "init chainsync error", err)
		logger.Fatal(err)
	}
	
	wg.Add(1)
	go cs.Run()
}
