package users

import (
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/lib/pq"
	"log"
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
		log.Print("failed to connect to postgres database")
		return nil, err
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Print("failed to start driver")
		return nil, err
	}
	m, err := migrate.NewWithDatabaseInstance("file://./migrations", "postgres", driver)
	if err != nil {
		log.Print("failed to start migrations")
		return nil, err
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Print("failed to run migrations")
		return nil, err
	}
	return &Postgres{db: db}, nil
}

func (p *Postgres) GetOrCreateUser(user *ExternalUser) (string, error) {
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
		log.Print("failed to get or create user")
		return "", err
	}
	return strconv.Itoa(userid), nil
}

func (p *Postgres) GetExternalUser(userID string) (*ExternalUser, error) {
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
		log.Print("failed to get external user")
		return nil, err
	}
	return ext, nil
}

func (p *Postgres) GetOrCreateGroup(group *ExternalGroup) (string, error) {
	query := `
	INSERT INTO groups (external_id, source)
	VALUES ($1, $2)
	ON CONFLICT(external_id, source)
		DO UPDATE SET source = $2
	RETURNING id;
	`
	var groupid int
	err := p.db.QueryRow(query, group.ID, int(group.Source)).Scan(&groupid)
	if err != nil {
		log.Print("failed to get or create group")
		return "", err
	}
	return strconv.Itoa(groupid), nil
}

func (p *Postgres) GetExternalGroup(groupID string) (*ExternalGroup, error) {
	query := `
	SELECT external_id, source FROM groups WHERE id = $1
	`
	ext := &ExternalGroup{}
	err := p.db.QueryRow(query, groupID).Scan(&ext.ID, &ext.Source)
	if err != nil {
		if err, ok := err.(*pq.Error); ok && err.Routine == "pg_atoi" {
			return nil, GroupNotFound
		}
		if err == sql.ErrNoRows {
			return nil, GroupNotFound
		}
		log.Print("failed to get external user")
		return nil, err
	}
	return ext, nil
}
