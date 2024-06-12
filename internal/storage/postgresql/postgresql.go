package postgresql

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type StoragePostgre struct {
	db *sql.DB
}

func New(storagePath string) (*StoragePostgre, error) {
	const op = "storage.postgresql.New"

	db, err := sql.Open("postgres", storagePath)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &StoragePostgre{db: db}, nil
}
