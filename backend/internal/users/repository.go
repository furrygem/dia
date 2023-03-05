package users

import "context"

type Reader interface {
	GetUserByID(id string, ctx context.Context) (*User, error)
	GetUserByUsername(username string, ctx context.Context) (*User, error)
	UsernameExists(username string, ctx context.Context) (bool, error)
}

type Writer interface {
	InsertUser(username string, hashedPassword string, ctx context.Context) (*User, error)
	SetActive(id string, active bool, ctx context.Context) (bool, error)
	UpdatePassword(id string, hashedPassword string, ctx context.Context) (string, error)
}

type Repository interface {
	Reader
	Writer
}
