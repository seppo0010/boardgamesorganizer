package users

import "errors"

type Source int

const (
	SourceCustom = iota
	SourceTelegram
)

var UserNotFound = errors.New("User not found")

type ExternalUser struct {
	ID     string
	Source Source
}

type Factory interface {
	GetOrCreate(user *ExternalUser) (string, error)
	GetExternal(userID string) (*ExternalUser, error)
}
