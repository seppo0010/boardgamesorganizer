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
