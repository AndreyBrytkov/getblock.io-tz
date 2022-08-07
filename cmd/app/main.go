package main

import (
	"sync"

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

}
