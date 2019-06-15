package meetings

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"os"
	"testing"
)

func getPostgres(t *testing.T) *Factory {
	URL := "postgres://testuser:testpass@localhost/testdb?sslmode=disable"
	if OS_URL := os.Getenv("BGO_TEST_POSTGRES_URL"); OS_URL != "" {
		URL = OS_URL
	}

	db, err := sql.Open("postgres", URL)
	if err != nil {
		t.Fatalf("cannot connect to db: %#v", err)
	}
	for _, table := range []string{"meetings", "schema_migrations", "attendees"} {
		_, err = db.Exec(fmt.Sprintf("DROP TABLE %s", table))
		if err, ok := err.(*pq.Error); ok && err.Routine != "DropErrorMsgNonExistent" {
			t.Fatalf("cannot drop db table %s: %#v", table, err)
		}
	}
	f, err := NewPostgres(&PostgresConfig{URL: URL, MigrationsPath: "./migrations"})
	if err != nil {
		t.Fatalf("cannot connect to db: %#v", err)
	}
	return f
}

func TestCreateGetDeleteMeetingPostgres(t *testing.T) {
	testCreateGetDeleteMeeting(t, getPostgres(t))
}

func TestAddRemoveAttendeePostgres(t *testing.T) {
	testAddRemoveAttendee(t, getPostgres(t))
}

func TestAttendeesPostgres(t *testing.T) {
	testAttendees(t, getPostgres(t))
}

func TestMeetingAlreadyActivePostgres(t *testing.T) {
	testMeetingAlreadyActive(t, getPostgres(t))
}

func TestMeetingInThePastPostgres(t *testing.T) {
	testMeetingInThePast(t, getPostgres(t))
}

func TestAddUserToMeetingBeforeMeetingPostgres(t *testing.T) {
	testMeetingInThePast(t, getPostgres(t))
}

func TestCannotAddAfterCapacityPostgres(t *testing.T) {
	testCannotAddAfterCapacity(t, getPostgres(t))
}

func TestMeetingIsClosedAfterStartPostgres(t *testing.T) {
	testMeetingIsClosedAfterStart(t, getPostgres(t))
}

func TestMeetingCannotRSVPAfterStartPostgres(t *testing.T) {
	testMeetingCannotRSVPAfterStart(t, getPostgres(t))
}

func TestCreateMeetingAfterClosedPostgres(t *testing.T) {
	testCreateMeetingAfterClosed(t, getPostgres(t))
}

func TestHaveMultipleClosedMeetingsPostgres(t *testing.T) {
	testHaveMultipleClosedMeetings(t, getPostgres(t))
}

func TestMeetingAttendeesDataPostgres(t *testing.T) {
	testMeetingAttendeesData(t, getPostgres(t))
}
