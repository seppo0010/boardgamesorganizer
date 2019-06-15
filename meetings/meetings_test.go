package meetings

import (
	ftime "github.com/seppo0010/boardgamesorganizer/time"
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
	"time"
)

func setTimeFactory(f *Factory) *ftime.Fake {
	tf := &ftime.Fake{CurrentNow: time.Date(2019, 5, 1, 17, 3, 7, 0, time.UTC)}
	f.SetTimeFactory(tf)
	return tf
}

func testCreateGetDeleteMeeting(t *testing.T, f *Factory) {
	assert := assert.New(t)
	setTimeFactory(f)

	m := &Meeting{
		Time: time.Date(2019, 5, 2, 20, 3, 7, 0, time.UTC),
	}
	groupID := "ashf"

	_, err := f.GetMeeting(groupID)
	assert.Equal(err, NoActiveMeeting)

	err = f.CreateMeeting(groupID, m)
	assert.NoError(err)

	m2, err := f.GetMeeting(groupID)
	assert.NoError(err)
	assert.Equal(m, m2)

	err = f.DeleteMeeting(groupID)
	assert.NoError(err)

	_, err = f.GetMeeting(groupID)
	assert.Equal(err, NoActiveMeeting)
}

func testAddRemoveAttendee(t *testing.T, f *Factory) {
	assert := assert.New(t)
	setTimeFactory(f)
	m := &Meeting{
		Time: time.Date(2019, 5, 2, 20, 3, 7, 0, time.UTC),
	}
	groupID := "ashf"
	userID := "oihf"

	err := f.CreateMeeting(groupID, m)
	assert.NoError(err)

	err = f.RemoveUserFromMeeting(groupID, userID)
	assert.Equal(err, UserDoesNotAttendMeeting)

	err = f.AddUserToMeeting(groupID, userID)
	assert.NoError(err)

	err = f.AddUserToMeeting(groupID, userID)
	assert.Equal(err, UserAlreadyAttendsMeeting)

	err = f.RemoveUserFromMeeting(groupID, userID)
	assert.NoError(err)

	err = f.RemoveUserFromMeeting(groupID, userID)
	assert.Equal(err, UserDoesNotAttendMeeting)
}

func testAttendees(t *testing.T, f *Factory) {
	assert := assert.New(t)
	setTimeFactory(f)
	m := &Meeting{
		Time: time.Date(2019, 5, 2, 20, 3, 7, 0, time.UTC),
	}
	groupID := "ashf"
	userID := "oihf"
	userID2 := "8126"

	err := f.CreateMeeting(groupID, m)
	assert.NoError(err)

	attendees, err := f.GetMeetingAttendees(groupID)
	assert.NoError(err)
	assert.Equal(attendees, []string{})

	err = f.AddUserToMeeting(groupID, userID)
	assert.NoError(err)

	attendees, err = f.GetMeetingAttendees(groupID)
	assert.NoError(err)
	assert.Equal(attendees, []string{userID})

	err = f.AddUserToMeeting(groupID, userID2)
	assert.NoError(err)

	attendees, err = f.GetMeetingAttendees(groupID)
	sort.Strings(attendees)
	expectedAttendees := []string{userID, userID2}
	sort.Strings(expectedAttendees)
	assert.NoError(err)
	assert.Equal(attendees, expectedAttendees)
}

func testMeetingAlreadyActive(t *testing.T, f *Factory) {
	assert := assert.New(t)
	setTimeFactory(f)
	m := &Meeting{
		Time: time.Date(2019, 5, 2, 20, 3, 7, 0, time.UTC),
	}
	groupID := "ashf"

	err := f.CreateMeeting(groupID, m)
	assert.NoError(err)

	err = f.CreateMeeting(groupID, m)
	assert.Equal(err, MeetingAlreadyActive)
}

func testMeetingInThePast(t *testing.T, f *Factory) {
	assert := assert.New(t)
	setTimeFactory(f)
	m := &Meeting{
		Time: time.Date(2019, 4, 2, 20, 3, 7, 0, time.UTC),
	}
	groupID := "ashf"

	err := f.CreateMeeting(groupID, m)
	assert.Equal(err, MeetingIsInThePast)
}

func testAddUserToMeetingBeforeMeeting(t *testing.T, f *Factory) {
	assert := assert.New(t)
	setTimeFactory(f)
	groupID := "ashf"
	userID := "oihf"

	err := f.AddUserToMeeting(groupID, userID)
	assert.Equal(err, NoActiveMeeting)
}

func testCannotAddAfterCapacity(t *testing.T, f *Factory) {
	assert := assert.New(t)
	setTimeFactory(f)
	m := &Meeting{
		Time:     time.Date(2019, 5, 2, 20, 3, 7, 0, time.UTC),
		Capacity: 2,
	}
	groupID := "ashf"
	userID := "qw"
	userID1 := "er"
	userID2 := "tr"

	err := f.CreateMeeting(groupID, m)
	assert.NoError(err)

	err = f.AddUserToMeeting(groupID, userID)
	assert.NoError(err)

	err = f.AddUserToMeeting(groupID, userID1)
	assert.NoError(err)

	err = f.AddUserToMeeting(groupID, userID2)
	assert.Equal(err, MeetingIsFull)
}

func testMeetingIsClosedAfterStart(t *testing.T, f *Factory) {
	assert := assert.New(t)
	tf := setTimeFactory(f)
	m := &Meeting{
		Time:     time.Date(2019, 5, 2, 20, 3, 7, 0, time.UTC),
		Capacity: 2,
	}
	groupID := "ashf"

	err := f.CreateMeeting(groupID, m)
	assert.NoError(err)

	tf.CurrentNow = time.Date(2019, 5, 3, 20, 3, 7, 0, time.UTC)

	_, err = f.GetMeeting(groupID)
	assert.Equal(err, NoActiveMeeting)
}

func testMeetingCannotRSVPAfterStart(t *testing.T, f *Factory) {
	assert := assert.New(t)
	tf := setTimeFactory(f)
	m := &Meeting{
		Time:     time.Date(2019, 5, 2, 20, 3, 7, 0, time.UTC),
		Capacity: 2,
	}
	groupID := "ashf"
	userID := "oihf"

	err := f.CreateMeeting(groupID, m)
	assert.NoError(err)

	tf.CurrentNow = time.Date(2019, 5, 3, 20, 3, 7, 0, time.UTC)

	err = f.AddUserToMeeting(groupID, userID)
	assert.Equal(err, NoActiveMeeting)

	err = f.AddUserToMeeting(groupID, userID)
	assert.Equal(err, NoActiveMeeting)
}

func testCreateMeetingAfterClosed(t *testing.T, f *Factory) {
	assert := assert.New(t)
	tf := setTimeFactory(f)
	m := &Meeting{
		Time: time.Date(2019, 5, 2, 20, 3, 7, 0, time.UTC),
	}
	groupID := "ashf"

	err := f.CreateMeeting(groupID, m)
	assert.NoError(err)

	tf.CurrentNow = time.Date(2019, 5, 3, 20, 3, 7, 0, time.UTC)

	m = &Meeting{
		Time: time.Date(2019, 5, 4, 20, 3, 7, 0, time.UTC),
	}
	err = f.CreateMeeting(groupID, m)
	assert.NoError(err)
}

func testHaveMultipleClosedMeetings(t *testing.T, f *Factory) {
	assert := assert.New(t)
	tf := setTimeFactory(f)
	groupID := "ashf"

	m := &Meeting{Time: time.Date(2019, 5, 2, 20, 3, 7, 0, time.UTC)}
	err := f.CreateMeeting(groupID, m)
	assert.NoError(err)

	tf.CurrentNow = time.Date(2019, 5, 3, 20, 3, 7, 0, time.UTC)

	m = &Meeting{Time: time.Date(2019, 5, 4, 20, 3, 7, 0, time.UTC)}
	err = f.CreateMeeting(groupID, m)
	assert.NoError(err)

	tf.CurrentNow = time.Date(2019, 5, 5, 20, 3, 7, 0, time.UTC)

	m = &Meeting{Time: time.Date(2019, 5, 6, 20, 3, 7, 0, time.UTC)}
	err = f.CreateMeeting(groupID, m)
	assert.NoError(err)
}

func testMeetingAttendeesData(t *testing.T, f *Factory) {
	assert := assert.New(t)
	setTimeFactory(f)
	groupID := "ashf"

	attendeesData := []string{"hello", "world"}
	err := f.SetMeetingAttendeesData(groupID, attendeesData)
	assert.Equal(err, NoActiveMeeting)

	m := &Meeting{Time: time.Date(2019, 5, 2, 20, 3, 7, 0, time.UTC)}
	err = f.CreateMeeting(groupID, m)
	assert.NoError(err)

	data := []string{}
	err = f.GetMeetingAttendeesData(groupID, &data)
	assert.NoError(err)
	assert.Equal(data, []string{})

	err = f.SetMeetingAttendeesData(groupID, attendeesData)
	assert.NoError(err)

	err = f.GetMeetingAttendeesData(groupID, &data)
	assert.NoError(err)
	assert.Equal(data, attendeesData)
}
