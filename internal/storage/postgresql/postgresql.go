package postgresql

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"time"
	"Ume/internal/storage"
)

type StoragePostgres struct {
	db *sql.DB
}

func NewPool(storagePath string) (*StoragePostgres, error) {
	const op = "storage.postgresql.New"

	db, err := sql.Open("postgres", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err_ping := db.Ping()
	if err_ping != nil {
		return nil, fmt.Errorf("%s: %w", op, err_ping)
	}

	return &StoragePostgres{db: db}, nil
}

func (s *StoragePostgres) AddUser(firstName, lastName, password, email, username string) (error) {
	const op = "storage.postgresql.AddUser"

	id, _ := s.GetUserIdByUsername(username)
	if id != 0 {
		fmt.Errorf("%s: username is already exists", op)
		return storage.ErrUserWithUsernameExists
	}

	stmt, err := s.db.Prepare(`INSERT INTO users (fisrt_name, last_name, email, password, date_birthday, username)
	 VALUES ($1, $2, $3, $4, $5, $6)`)
	if err != nil {
		return fmt.Errorf("%s: error with preparing query %w", op, err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(firstName, lastName, email, password, time.Now(), username)
	if err != nil {
		return fmt.Errorf("%s: error with adding new user %w", op, err)
	}
	
	return nil
}

func (s *StoragePostgres) GetUserIdByUsername(username string) (int64, error) {
	const op = "storage.postgresql.getUserIdByUsername"

	stmt, err := s.db.Prepare(`SELECT id FROM users WHERE username = $1::varchar`)
	if err != nil {
		return 0, fmt.Errorf("%s: error with preparing query %w", op, err)
	}
	defer stmt.Close()

	var id int64
	err = stmt.QueryRow(username).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: error with gettind user id %w", op, err)
	}

	return id, nil
}


func (s *StoragePostgres) AddFriend(userNameFirst, userNameSecond string) (error) {
	const op = "storage.postgresql.AddFriend"

	userIdFirst, err := s.GetUserIdByUsername(userNameFirst)
	if err != nil {
		return fmt.Errorf("%s: error with getting user id with username %s: %w", op, userNameFirst, err)
	}

	userIdSecond, err := s.GetUserIdByUsername(userNameSecond)
	if err != nil {
		return fmt.Errorf("%s: error with getting user id with username %s: %w", op, userNameSecond, err)
	}

	stmt, err := s.db.Prepare(`INSERT INTO friends (user_id_1, user_id_2) values ($1::int, $2::int)`)
	if err != nil {
		return fmt.Errorf("%s: error with preparing query %w", op, err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(userIdFirst, userIdSecond)
	if err != nil {
		return fmt.Errorf("%s: error with making friends for users %s and %s %w", op, userNameFirst, userNameSecond, err)
	}

	return nil
}

func (s *StoragePostgres) AddMessage(userNameFirst, userNameSecond, messageText string) (error) {
	const op = "storage.postgresql.AddFriend"

	userIdFirst, err := s.GetUserIdByUsername(userNameFirst)
	if err != nil {
		return fmt.Errorf("%s: error with getting user id with username %w", op, userNameFirst, err)
	}

	userIdSecond, err := s.GetUserIdByUsername(userNameSecond)
	if err != nil {
		return fmt.Errorf("%s: error with getting user id with username %w", op, userNameSecond, err)
	}

	stmt, err := s.db.Prepare(`INSERT INTO messages (user_id_to, user_id_from, text_message) values ($1::int, $2::int, $3::text)`)
	if err != nil {
		return fmt.Errorf("%s: error with preparing query %w", op, err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(userIdSecond, userIdFirst, messageText)
	if err != nil {
		return fmt.Errorf("%s: error with adding message for users %s and %s %w", op, userNameFirst, userNameSecond, err)
	}

	return nil
}

func (s *StoragePostgres) FindUserPassword(username string) (string, error) {
	const op = "storage.postgresql.FindUser"

	stmt, err := s.db.Prepare(`SELECT password FROM users WHERE username = $1`)
	if err != nil {
		return "", fmt.Errorf("%s: error with preparing query %w", op, err)
	}
	defer stmt.Close()

	var password string
	err = stmt.QueryRow(username).Scan(&password)
	if err != nil {
		if err == sql.ErrNoRows{
			fmt.Errorf("%s: not found user with username %s or incorrect password %w", op, username, err)
			return "", storage.ErrUserNotExist
		}
		return "", fmt.Errorf("%s: error with gettind user %w", op, err)
	}

	return password, nil
}