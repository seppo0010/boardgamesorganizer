package users

import (
	"fmt"
	"os"
	"testing"
)

func getPostgres(t *testing.T) *Postgres {
	URL := "postgres://testuser:testpass@localhost/testdb?sslmode=disable"
	if OS_URL := os.Getenv("BGO_TEST_POSTGRES_URL"); OS_URL != "" {
		URL = OS_URL
	}

	f, err := NewPostgres(&PostgresConfig{URL: URL})
	if err != nil {
		t.Fatalf("cannot connect to db: %#v", err)
	}
	for _, table := range []string{"users", "schema_migrations"} {
		_, err = f.db.Exec(fmt.Sprintf("DROP TABLE %s", table))
		if err != nil {
			t.Fatalf("cannot drop db users: %#v", err)
		}
	}
	f, err = NewPostgres(&PostgresConfig{URL: URL})
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
