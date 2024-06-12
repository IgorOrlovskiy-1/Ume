package storage

import "errors"

var (
	ErrUserWithLoginExists = errors.New("user with login already exists")
	ErrNotCorrectPassword  = errors.New("not correct password")
	ErrChatWithUserExists  = errors.New("chat with user already exists")
	ErrUserNotExist        = errors.New("user not exist")
	ErrSendMessage         = errors.New("failed to send message")
)
