package users

import (
	"fmt"
	"os"
	"testing"
	"database/sql"
)

func getPostgres(t *testing.T) *Postgres {
	URL := "postgres://testuser:testpass@localhost/testdb?sslmode=disable"
	if OS_URL := os.Getenv("BGO_TEST_POSTGRES_URL"); OS_URL != "" {
		URL = OS_URL
	}

	db, err := sql.Open("postgres", URL)
	if err != nil {
		t.Fatalf("cannot connect to db: %#v", err)
	}
	for _, table := range []string{"users", "schema_migrations", "groups"} {
		_, err = db.Exec(fmt.Sprintf("DROP TABLE %s", table))
		if err != nil {
			t.Fatalf("cannot drop db users: %#v", err)
		}
	}
	f, err := NewPostgres(&PostgresConfig{URL: URL})
	if err != nil {
		t.Fatalf("cannot connect to db: %#v", err)
	}
	return f
}

func TestCreateUserPostgres(t *testing.T) {
	testCreateUser(t, getPostgres(t))
}

func TestCreateGetUserPostgres(t *testing.T) {
	testCreateGetUser(t, getPostgres(t))
}

func TestCreateSecondUserPostgres(t *testing.T) {
	testCreateSecondUser(t, getPostgres(t))
}

func TestCreateSecondUserSourcePostgres(t *testing.T) {
	testCreateSecondUserSource(t, getPostgres(t))
}

func TestGetExistingUserPostgres(t *testing.T) {
	testGetExistingUser(t, getPostgres(t))
}

func TestGetNoExistingUserPostgres(t *testing.T) {
	testGetNoExistingUser(t, getPostgres(t))
}

func TestCreateGroupPostgres(t *testing.T) {
	testCreateGroup(t, getPostgres(t))
}

func TestCreateGetGroupPostgres(t *testing.T) {
	testCreateGetGroup(t, getPostgres(t))
}

func TestCreateSecondGroupPostgres(t *testing.T) {
	testCreateSecondGroup(t, getPostgres(t))
}

func TestCreateSecondGroupSourcePostgres(t *testing.T) {
	testCreateSecondGroupSource(t, getPostgres(t))
}

func TestGetExistingGroupPostgres(t *testing.T) {
	testGetExistingGroup(t, getPostgres(t))
}

func TestGetNoExistingGroupPostgres(t *testing.T) {
	testGetNoExistingGroup(t, getPostgres(t))
}
