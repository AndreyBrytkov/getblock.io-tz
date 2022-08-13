package main

import (
	"sync"

	"github.com/AndreyBrytkov/getblock.io-tz/internal/chainsync"
	"github.com/AndreyBrytkov/getblock.io-tz/internal/repository"
	"github.com/AndreyBrytkov/getblock.io-tz/internal/usecase"
	"github.com/AndreyBrytkov/getblock.io-tz/pkg/utils"
	"github.com/AndreyBrytkov/getblock.io-tz/internal/restserver"
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

	uc := usecase.GetUsecase(logger, &config.AppConfig, repo.Api, repo.Storage)

	server := restserver.NewRestApi(logger, &config.ServerConfig, wg, uc)

	cs := chainsync.GetChainSynchronizer(logger, &config.AppConfig, wg, repo.Api, repo.Storage)

	err = cs.Init()
	if err != nil {
		err = utils.WrapErr(caller, "init chainsync error", err)
		logger.Fatal(err)
	}

	wg.Add(2)
	go cs.Run()
	go server.Run()
}
