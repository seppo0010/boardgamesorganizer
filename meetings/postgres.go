package meetings

import (
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"time"
)

type PostgresConfig struct {
	URL string
}

type Postgres struct {
	db *sql.DB
}

func NewPostgres(config *PostgresConfig) (*Factory, error) {
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
	return NewFactory(&Postgres{db: db}), nil
}

func (p *Postgres) CreateMeeting(groupID string, meeting *Meeting) error {
	query := `
	INSERT INTO meetings (group_id, time, location)
	VALUES ($1, $2, $3)
	ON CONFLICT DO NOTHING
	RETURNING id;
	`
	id := 0
	err := p.db.QueryRow(query, groupID, meeting.Time, meeting.Location).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return MeetingAlreadyActive
		}
		log.Printf("failed to create meeting: %#v", err)
		return UnexpectedError
	}
	return nil
}

func (p *Postgres) DeleteMeeting(groupID string) error {
	query := `
	DELETE FROM meetings WHERE group_id = $1
	`
	result, err := p.db.Exec(query, groupID)
	if err != nil {
		log.Printf("failed to delete meeting: %#v", err)
		return UnexpectedError
	}
	affectedRows, err := result.RowsAffected()
	if err != nil {
		log.Printf("failed to get affected rows after deleting meeting: %#v", err)
		return UnexpectedError
	}
	if affectedRows == 0 {
		return NoActiveMeeting
	}
	return nil
}

func (p *Postgres) GetMeeting(groupID string) (*Meeting, error) {
	query := `
	SELECT time AT TIME ZONE 'GMT', location FROM meetings WHERE group_id = $1
	`
	m := &Meeting{}
	err := p.db.QueryRow(query, groupID).Scan(&m.Time, &m.Location)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NoActiveMeeting
		}
		log.Printf("failed to get meeting: %#v", err)
		return nil, UnexpectedError
	}
	m.Time = m.Time.In(time.UTC)
	return m, nil
}

func (p *Postgres) AddUserToMeeting(groupID string, userID string) error {
	query := `
	INSERT INTO attendees (group_id, user_id)
	VALUES ($1, $2)
	ON CONFLICT DO NOTHING
	RETURNING id;
	`
	id := 0
	err := p.db.QueryRow(query, groupID, userID).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return UserAlreadyAttendsMeeting
		}
		log.Printf("failed to add attendee: %#v", err)
		return UnexpectedError
	}
	return nil
}
func (p *Postgres) RemoveUserFromMeeting(groupID string, userID string) error {
	query := `
	DELETE FROM attendees WHERE group_id = $1 AND user_id = $2
	`
	result, err := p.db.Exec(query, groupID, userID)
	if err != nil {
		log.Printf("failed to delete attendee: %#v", err)
		return UnexpectedError
	}
	affectedRows, err := result.RowsAffected()
	if err != nil {
		log.Printf("failed to get affected rows after deleting attendee: %#v", err)
		return UnexpectedError
	}
	if affectedRows == 0 {
		return UserDoesNotAttendMeeting
	}
	return nil
}

func (p *Postgres) GetMeetingAttendees(groupID string) ([]string, error) {
	query := `
	SELECT user_id FROM attendees WHERE group_id = $1
	`
	rows, err := p.db.Query(query, groupID)
	if err != nil {
		log.Printf("failed to get attendees: %#v", err)
		return nil, UnexpectedError
	}
	defer rows.Close()
	userIDs := make([]string, 0)
	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			log.Printf("failed to get next attendee: %#v", err)
			return nil, UnexpectedError
		}
		userIDs = append(userIDs, userID)
	}
	if err := rows.Err(); err != nil {
		log.Printf("failed to get close attendees: %#v", err)
		return nil, UnexpectedError
	}
	return userIDs, nil

}
