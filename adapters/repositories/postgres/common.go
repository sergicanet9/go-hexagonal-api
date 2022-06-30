package postgres

import "database/sql"

// PostgresRepository struct of a mongo repository
type PostgresRepository struct {
	db *sql.DB
}
