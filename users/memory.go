package users

import "strconv"

type Memory struct {
	users map[string]*ExternalUser
}

func NewMemory() *Memory {
	return &Memory{users: map[string]*ExternalUser{}}
}

func (m *Memory) GetOrCreate(user *ExternalUser) (string, error) {
	for id, aUser := range m.users {
		if user.ID == aUser.ID && user.Source == aUser.Source {
			return id, nil
		}
	}
	ID := strconv.Itoa(len(m.users) + 1)
	m.users[ID] = user
	return ID, nil
}

func (m *Memory) GetExternal(userID string) (*ExternalUser, error) {
	user, found := m.users[userID]
	if !found {
		return nil, UserNotFound
	}
	return user, nil
}
