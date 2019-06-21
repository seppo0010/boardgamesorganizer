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
	Closed   bool
}

type Attendee struct {
	UserID string
	Amount int
}

type Inner interface {
	CreateMeeting(groupID string, meeting *Meeting) error
	DeleteMeeting(groupID string) error
	GetMeeting(groupID string) (*Meeting, error)
	SetMeetingAttendeesData(groupID string, data interface{}) error
	GetMeetingAttendeesData(groupID string, data interface{}) error
	UserRSVPMeeting(groupID string, attendee *Attendee) error
	GetMeetingAttendees(groupID string) ([]*Attendee, error)
	CloseMeeting(groupID string) error
}

type Factory struct {
	Inner
	timeFactory ftime.Factory
}

func NewFactory(inner Inner) *Factory {
	return &Factory{Inner: inner, timeFactory: ftime.NewReal()}
}

func (f *Factory) GetMeeting(groupID string) (*Meeting, error) {
	return f.closeMeetingIfNeeded(groupID)
}

func (f *Factory) closeMeetingIfNeeded(groupID string) (*Meeting, error) {
	meeting, err := f.Inner.GetMeeting(groupID)
	if err != nil {
		return nil, err
	}
	if !meeting.Closed && meeting.Time.Before(f.timeFactory.Now()) {
		err = f.CloseMeeting(groupID)
		if err != nil {
			return nil, err
		}
		return nil, NoActiveMeeting
	}
	return meeting, nil
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
	if _, err := f.closeMeetingIfNeeded(groupID); err != nil && err != NoActiveMeeting {
		return err
	}
	if meeting.Time.Before(f.timeFactory.Now()) {
		return MeetingIsInThePast
	}
	return f.Inner.CreateMeeting(groupID, meeting)
}

func (f *Factory) UserRSVPMeeting(groupID string, attendee *Attendee) error {
	meeting, err := f.GetMeeting(groupID)
	if err != nil {
		return err
	}
	attendees, err := f.GetMeetingAttendees(groupID)
	if err != nil {
		return err
	}
	taken := 0
	for _, att := range attendees {
		if att.UserID != attendee.UserID {
			taken += att.Amount
		}
	}
	if meeting.Capacity > 0 && meeting.Capacity < taken+attendee.Amount {
		return MeetingIsFull
	}
	// FIXME: possible race condition if two people RSVP at the same time
	return f.Inner.UserRSVPMeeting(groupID, attendee)
}
