package postgres

import (
	"fmt"

	"github.com/AndreyBrytkov/getblock.io-tz/internal/adapter"
	"github.com/AndreyBrytkov/getblock.io-tz/internal/models"
	"github.com/AndreyBrytkov/getblock.io-tz/pkg/utils"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const caller = "postgres"

// Main package object
type Postgres struct {
	logger *utils.MyLogger
	config *models.StorageConfig
	DB     *sqlx.DB
}

// Getter
func GetPostgres(l *utils.MyLogger, c *models.StorageConfig) adapter.Storage {
	// connect db
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", c.Address, c.Port, c.User, c.Password, c.DB)
	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		l.Fatal(utils.WrapErr(caller, "connect db error", err))
		return nil
	}
	l.Info(caller, "connection established...")

	// create tables
	schemas := []string{SCHEMA_TRANSACTIONS}
	for _, schema := range schemas {
		_, err := db.Exec(schema)
		if err != nil {
			l.Fatal(utils.WrapErr(caller, "create table error", err))
			return nil
		}
	}

	return &Postgres{
		logger: l,
		config: c,
		DB:     db,
	}
}
