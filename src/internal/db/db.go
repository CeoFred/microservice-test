package db

import (
	"database/sql"

	"github.com/SAMBA-Research/microservice-template/internal/config"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func NewDbConnection(cfg *config.Config) (db *bun.DB) {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(cfg.DatabaseDSN)))

	db = bun.NewDB(sqldb, pgdialect.New())
	return
}
