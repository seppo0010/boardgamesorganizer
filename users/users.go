package users

import "errors"

type Source int

const (
	SourceCustom = iota
	SourceTelegram
)

var UserNotFound = errors.New("User not found")
var GroupNotFound = errors.New("Group not found")

type ExternalUser struct {
	ID          string
	Source      Source
	DisplayName string
}

type ExternalGroup struct {
	ID     string
	Source Source
}

type Factory interface {
	GetOrCreateUser(user *ExternalUser) (string, error)
	GetExternalUser(userID string) (*ExternalUser, error)
	GetOrCreateGroup(group *ExternalGroup) (string, error)
	GetExternalGroup(groupID string) (*ExternalGroup, error)
	GetUsers(userIDs []string) (map[string]*ExternalUser, error)
}
