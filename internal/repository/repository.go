package repository

import (
	"github.com/AndreyBrytkov/getblock.io-tz/internal/adapter"
	"github.com/AndreyBrytkov/getblock.io-tz/internal/models"
	"github.com/AndreyBrytkov/getblock.io-tz/internal/repository/getblockapi"
	"github.com/AndreyBrytkov/getblock.io-tz/internal/repository/postgres"
	"github.com/AndreyBrytkov/getblock.io-tz/pkg/utils"
)

const caller = "repository"

// Main package object
type Repository struct {
	Storage adapter.Storage
	Api     adapter.GetBlockApi
}

//Getter
func GetRepository(logger *utils.MyLogger, config *models.Config) *Repository {
	var storage adapter.Storage
	switch config.StorageConfig.Type {
	default:
		storage = postgres.GetPostgres(logger, &config.StorageConfig)
	}

	return &Repository{
		Storage: storage,
		Api:     getblockapi.GetGetBlockApi(logger, &config.GetBlockConfig),
	}
}
