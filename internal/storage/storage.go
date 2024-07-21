package storage

import "errors"

var (
	ErrUserWithUsernameExists = errors.New("user with username already exists")
	ErrChatWithUserExists  = errors.New("chat with user already exists")
	ErrUserNotExist        = errors.New("user not exist")
	ErrSendMessage         = errors.New("failed to send message")
)
