package users

import (
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/lib/pq"
	"strconv"
)

type PostgresConfig struct {
	URL string
}

type Postgres struct {
	db *sql.DB
}

func NewPostgres(config *PostgresConfig) (*Postgres, error) {
	db, err := sql.Open("postgres", config.URL)
	if err != nil {
		return nil, err
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, err
	}
	m, err := migrate.NewWithDatabaseInstance("file://./migrations", "postgres", driver)
	if err != nil {
		return nil, err
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return nil, err
	}
	return &Postgres{db: db}, nil
}

func (p *Postgres) GetOrCreate(user *ExternalUser) (string, error) {
	query := `
	INSERT INTO users (external_id, source)
	VALUES ($1, $2)
	ON CONFLICT(external_id, source)
		DO UPDATE SET source = $2
	RETURNING id;
	`
	var userid int
	err := p.db.QueryRow(query, user.ID, int(user.Source)).Scan(&userid)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(userid), nil
}

func (p *Postgres) GetExternal(userID string) (*ExternalUser, error) {
	query := `
	SELECT external_id, source FROM users WHERE id = $1
	`
	ext := &ExternalUser{}
	err := p.db.QueryRow(query, userID).Scan(&ext.ID, &ext.Source)
	if err != nil {
		if err, ok := err.(*pq.Error); ok && err.Routine == "pg_atoi" {
			return nil, UserNotFound
		}
		if err == sql.ErrNoRows {
			return nil, UserNotFound
		}
		return nil, err
	}
	return ext, nil
}
