package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Config struct {
	Uri string
}

// Open the sqlite db
func Open(cfg Config) (*sqlx.DB, error) {
	return sqlx.Open("sqlite3", cfg.Uri)
}
