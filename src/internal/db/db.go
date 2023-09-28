package db

import (
	"context"
	"database/sql"
	"github.com/SAMBA-Research/microservice-template/internal/config"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
)

func NewDbConnection(cfg *config.Config) (db *bun.DB) {
	sqldb, err := sql.Open(sqliteshim.ShimName, "file::memory:?cache=shared")
	if err != nil {
		panic(err)
	}
	ctx := context.Background()

	db = bun.NewDB(sqldb, sqlitedialect.New())

	_, err = db.NewCreateTable().Model((*Message)(nil)).Exec(ctx)

	if err != nil {
		panic(err)
	}
	return
}
