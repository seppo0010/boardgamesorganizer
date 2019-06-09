package meetings

type Memory struct {
	groupMeetings  map[string]*Meeting
	groupAttendees map[string][]string
}

func stringInSliceIndex(a string, list []string) int {
	for i, b := range list {
		if b == a {
			return i
		}
	}
	return -1
}

func NewMemory() *Factory {
	return NewFactory(&Memory{
		groupMeetings:  map[string]*Meeting{},
		groupAttendees: map[string][]string{},
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

func (m *Memory) AddUserToMeeting(groupID string, userID string) error {
	if _, found := m.groupMeetings[groupID]; !found {
		return NoActiveMeeting
	}
	if attendees, found := m.groupAttendees[groupID]; found {
		if stringInSliceIndex(userID, attendees) != -1 {
			return UserAlreadyAttendsMeeting
		}
	} else {
		m.groupAttendees[groupID] = []string{}
	}
	m.groupAttendees[groupID] = append(m.groupAttendees[groupID], userID)
	return nil
}
func (m *Memory) RemoveUserFromMeeting(groupID string, userID string) error {
	if _, found := m.groupMeetings[groupID]; !found {
		return NoActiveMeeting
	}
	attendees, found := m.groupAttendees[groupID]
	if !found {
		return UserDoesNotAttendMeeting
	}
	index := stringInSliceIndex(userID, attendees)
	if index == -1 {
		return UserDoesNotAttendMeeting
	}
	attendees[index] = attendees[len(attendees)-1]
	m.groupAttendees[groupID] = attendees[:len(attendees)-1]
	return nil
}

func (m *Memory) GetMeetingAttendees(groupID string) ([]string, error) {
	if _, found := m.groupMeetings[groupID]; !found {
		return nil, NoActiveMeeting
	}
	attendees, found := m.groupAttendees[groupID]
	if !found {
		return []string{}, nil
	}
	return attendees, nil

}
