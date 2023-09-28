package db

import "github.com/uptrace/bun"

type Message struct {
	bun.BaseModel `bun:"table:messages,alias:m"`
	ID            int64  `bun:",pk,autoincrement"`
	Message       string `json:"message,omitempty"`
}
