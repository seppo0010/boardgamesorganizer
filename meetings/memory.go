package meetings

import (
	"encoding/json"
	"log"
)

type Memory struct {
	groupMeetings        map[string]*Meeting
	groupAttendees       map[string][]*Attendee
	closedMeetings       map[string][]*Meeting
	meetingAttendeesData map[string][]byte
}

func NewMemory() *Factory {
	return NewFactory(&Memory{
		groupMeetings:        map[string]*Meeting{},
		groupAttendees:       map[string][]*Attendee{},
		closedMeetings:       map[string][]*Meeting{},
		meetingAttendeesData: map[string][]byte{},
	})
}

func (m *Memory) CreateMeeting(groupID string, meeting *Meeting) error {
	if _, found := m.groupMeetings[groupID]; found {
		return MeetingAlreadyActive
	}
	m.groupMeetings[groupID] = meeting
	return nil
}

func (m *Memory) DeleteMeeting(groupID string) error {
	if _, found := m.groupMeetings[groupID]; !found {
		return NoActiveMeeting
	}
	delete(m.groupMeetings, groupID)
	return nil
}

func (m *Memory) GetMeeting(groupID string) (*Meeting, error) {
	meeting, found := m.groupMeetings[groupID]
	if !found {
		return nil, NoActiveMeeting
	}
	return meeting, nil
}

func (m *Memory) UserRSVPMeeting(groupID string, attendee *Attendee) error {
	if _, found := m.groupMeetings[groupID]; !found {
		return NoActiveMeeting
	}
	if attendees, found := m.groupAttendees[groupID]; found {
		for _, att := range attendees {
			if attendee.UserID == att.UserID {
				if attendee.Amount == att.Amount {
					if att.Amount == 0 {
						return UserDoesNotAttendMeeting
					}
					return UserAlreadyAttendsMeeting
				}
				att.Amount = attendee.Amount
				return nil
			}
		}
	} else {
		m.groupAttendees[groupID] = []*Attendee{}
	}
	if attendee.Amount == 0 {
		return UserDoesNotAttendMeeting
	}
	m.groupAttendees[groupID] = append(m.groupAttendees[groupID], attendee)
	return nil
}
func (m *Memory) GetMeetingAttendees(groupID string) ([]*Attendee, error) {
	if _, found := m.groupMeetings[groupID]; !found {
		return nil, NoActiveMeeting
	}
	attendees, found := m.groupAttendees[groupID]
	if !found {
		return []*Attendee{}, nil
	}
	return attendees, nil

}
func (m *Memory) CloseMeeting(groupID string) error {
	if _, found := m.groupMeetings[groupID]; !found {
		return NoActiveMeeting
	}
	if _, found := m.closedMeetings[groupID]; !found {
		m.closedMeetings[groupID] = []*Meeting{}
	}
	m.closedMeetings[groupID] = append(m.closedMeetings[groupID], m.groupMeetings[groupID])
	delete(m.groupMeetings, groupID)
	return nil
}
func (m *Memory) SetMeetingAttendeesData(groupID string, data interface{}) error {
	if _, found := m.groupMeetings[groupID]; !found {
		return NoActiveMeeting
	}
	v, err := json.Marshal(data)
	if err != nil {
		log.Print(err)
		return err
	}
	m.meetingAttendeesData[groupID] = v
	return nil
}
func (m *Memory) GetMeetingAttendeesData(groupID string, v interface{}) error {
	if _, found := m.groupMeetings[groupID]; !found {
		return NoActiveMeeting
	}
	data, found := m.meetingAttendeesData[groupID]
	if !found {
		return nil
	}
	err := json.Unmarshal(data, v)
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}
