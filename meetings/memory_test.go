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

func TestCannotAddAfterCapacityMemory(t *testing.T) {
	testCannotAddAfterCapacity(t, NewMemory())
}

func TestMeetingIsClosedAfterStartMemory(t *testing.T) {
	testMeetingIsClosedAfterStart(t, NewMemory())
}

func TestMeetingCannotRSVPAfterStartMemory(t *testing.T) {
	testMeetingCannotRSVPAfterStart(t, NewMemory())
}

func TestCreateMeetingAfterClosedMemory(t *testing.T) {
	testCreateMeetingAfterClosed(t, NewMemory())
}

func TestHaveMultipleClosedMeetingsMemory(t *testing.T) {
	testHaveMultipleClosedMeetings(t, NewMemory())
}
