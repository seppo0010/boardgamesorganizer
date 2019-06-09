package meetings

import (
	"testing"
)

func TestCreateGetDeleteMeetingMemory(t *testing.T) {
	testCreateGetDeleteMeeting(t, NewMemory())
}

func TestAddRemoveAttendeeMemory(t *testing.T) {
	testAddRemoveAttendee(t, NewMemory())
}

func TestAttendeesMemory(t *testing.T) {
	testAttendees(t, NewMemory())
}

func TestMeetingAlreadyActiveMemory(t *testing.T) {
	testMeetingAlreadyActive(t, NewMemory())
}

func TestMeetingInThePastMemory(t *testing.T) {
	testMeetingInThePast(t, NewMemory())
}

func TestAddUserToMeetingBeforeMeetingMemory(t *testing.T) {
	testMeetingInThePast(t, NewMemory())
}
