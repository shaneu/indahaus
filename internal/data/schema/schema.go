package schema

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

const schema = `
	CREATE TABLE ip_results (
		ip_address TEXT PRIMARY KEY,
		id TEXT UNIQUE,
		created_at DATETIME,
		updated_at DATETIME,
		response_codes TEXT
	)
`

// Migrate creates our table
func Migrate(db *sqlx.DB) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf("unable to migrate db %v", r)
		}
	}()

	db.MustExec(schema)

	return nil
}
