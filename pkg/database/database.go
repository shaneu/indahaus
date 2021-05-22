package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Config struct {
	Uri      string
	Username string
	Password string
}

func Open(cfg Config) (*sqlx.DB, error) {
	u := fmt.Sprintf("%s&_auth_user=%s&_auth_pass=%s", cfg.Uri, cfg.Username, cfg.Password)
	return sqlx.Open("sqlite3", u)
}
