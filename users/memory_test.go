package users

import (
	"testing"
)

func TestCreateUserMemory(t *testing.T) {
	t.Parallel()
	testCreateUser(t, NewMemory())
}

func TestCreateGetUserMemory(t *testing.T) {
	t.Parallel()
	testCreateGetUser(t, NewMemory())
}

func TestCreateSecondUserMemory(t *testing.T) {
	t.Parallel()
	testCreateSecondUser(t, NewMemory())
}

func TestCreateSecondUserSourceMemory(t *testing.T) {
	t.Parallel()
	testCreateSecondUserSource(t, NewMemory())
}

func TestGetExistingUserMemory(t *testing.T) {
	t.Parallel()
	testGetExistingUser(t, NewMemory())
}

func TestGetNoExistingUserMemory(t *testing.T) {
	t.Parallel()
	testGetNoExistingUser(t, NewMemory())
}

func TestCreateGroupMemory(t *testing.T) {
	t.Parallel()
	testCreateGroup(t, NewMemory())
}

func TestCreateGetGroupMemory(t *testing.T) {
	t.Parallel()
	testCreateGetGroup(t, NewMemory())
}

func TestCreateSecondGroupMemory(t *testing.T) {
	t.Parallel()
	testCreateSecondGroup(t, NewMemory())
}

func TestCreateSecondGroupSourceMemory(t *testing.T) {
	t.Parallel()
	testCreateSecondGroupSource(t, NewMemory())
}

func TestGetExistingGroupMemory(t *testing.T) {
	t.Parallel()
	testGetExistingGroup(t, NewMemory())
}

func TestGetNoExistingGroupMemory(t *testing.T) {
	t.Parallel()
	testGetNoExistingGroup(t, NewMemory())
}

func TestGetUsersMemory(t *testing.T) {
	t.Parallel()
	testGetUsers(t, NewMemory())
}
