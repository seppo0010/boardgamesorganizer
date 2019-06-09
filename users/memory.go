package users

import "strconv"

type Memory struct {
	users  map[string]*ExternalUser
	groups map[string]*ExternalGroup
}

func NewMemory() *Memory {
	return &Memory{
		users:  map[string]*ExternalUser{},
		groups: map[string]*ExternalGroup{},
	}
}

func (m *Memory) GetOrCreateUser(user *ExternalUser) (string, error) {
	for id, aUser := range m.users {
		if user.ID == aUser.ID && user.Source == aUser.Source {
			return id, nil
		}
	}
	ID := strconv.Itoa(len(m.users) + 1)
	m.users[ID] = user
	return ID, nil
}

func (m *Memory) GetExternalUser(userID string) (*ExternalUser, error) {
	user, found := m.users[userID]
	if !found {
		return nil, UserNotFound
	}
	return user, nil
}

func (m *Memory) GetOrCreateGroup(group *ExternalGroup) (string, error) {
	for id, aGroup := range m.groups {
		if group.ID == aGroup.ID && group.Source == aGroup.Source {
			return id, nil
		}
	}
	ID := strconv.Itoa(len(m.groups) + 1)
	m.groups[ID] = group
	return ID, nil
}

func (m *Memory) GetExternalGroup(groupID string) (*ExternalGroup, error) {
	group, found := m.groups[groupID]
	if !found {
		return nil, GroupNotFound
	}
	return group, nil
}
