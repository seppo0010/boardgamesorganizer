package meetings

import (
	"errors"
	ftime "github.com/seppo0010/boardgamesorganizer/time"
	"time"
)

var NoActiveMeeting = errors.New("Group has no active meeting")
var MeetingAlreadyActive = errors.New("Group can only have one active meeting at a time")
var MeetingIsInThePast = errors.New("Meetings can only be created in the future")
var UserAlreadyAttendsMeeting = errors.New("User is already attending meeting")
var UserDoesNotAttendMeeting = errors.New("User is not attending meeting")
var UnexpectedError = errors.New("Unexpected error")
var MeetingIsFull = errors.New("Meeting is full")

type Meeting struct {
	Time     time.Time
	Location string
	Capacity int
}

type Inner interface {
	CreateMeeting(groupID string, meeting *Meeting) error
	DeleteMeeting(groupID string) error
	GetMeeting(groupID string) (*Meeting, error)
	AddUserToMeeting(groupID string, userID string) error
	RemoveUserFromMeeting(groupID string, userID string) error
	GetMeetingAttendees(groupID string) ([]string, error)
}

type Factory struct {
	Inner
	timeFactory ftime.Factory
}

func NewFactory(inner Inner) *Factory {
	return &Factory{Inner: inner, timeFactory: ftime.NewReal()}
}

func (f *Factory) SetTimeFactory(tf ftime.Factory) {
	f.timeFactory = tf
}

func (f *Factory) CanCreateMeeting(groupID string, meeting *Meeting) error {
	if meeting.Time.Before(f.timeFactory.Now()) {
		return MeetingIsInThePast
	}
	if _, err := f.GetMeeting(groupID); err == nil {
		return MeetingAlreadyActive
	}
	return nil
}
func (f *Factory) CreateMeeting(groupID string, meeting *Meeting) error {
	if meeting.Time.Before(f.timeFactory.Now()) {
		return MeetingIsInThePast
	}
	return f.Inner.CreateMeeting(groupID, meeting)
}

func (f *Factory) AddUserToMeeting(groupID string, userID string) error {
	meeting, err := f.GetMeeting(groupID)
	if err != nil {
		return err
	}
	attendees, err := f.GetMeetingAttendees(groupID)
	if err != nil {
		return err
	}
	if meeting.Capacity > 0 && meeting.Capacity <= len(attendees) {
		return MeetingIsFull
	}
	return f.Inner.AddUserToMeeting(groupID, userID)
}
