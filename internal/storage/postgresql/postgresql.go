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

	err_ping := db.Ping()

	if err_ping != nil {
		return nil, fmt.Errorf("%s: %w", op, err_ping)
	}

	return &StoragePostgre{db: db}, nil
}

func (s *StoragePostgre) AddUser(firstName, lastName, password string) (int64, error) {
	const op = "storage.postgresql.AddUser"

	stmt, err := s.db.Prepare(`INSERT INTO main.users (first_name, last_name, password) VALUES ($1, $2, $3)`)

	if err != nil {
		return 0, fmt.Errorf("%s: error with preparing query %w", op, err)
	}

	res, err := stmt.Exec(firstName, lastName, password)
	if err != nil {
		return 0, fmt.Errorf("%s: error with adding new user %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: error with getting id of last raw %w", op, err)
	}

	return id, nil
}
