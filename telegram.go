package main

import (
	"errors"
	"fmt"
	"github.com/seppo0010/boardgamesorganizer/meetings"
	"github.com/seppo0010/boardgamesorganizer/users"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"strconv"
	"strings"
	"time"
)

var ErrNeedsSegments = errors.New("Needs to be location;datetime;capacity. For example 'Home;2019-03-05 20:01:00;8'")
var ErrInvalidDate = errors.New("Datetime must follow the format YYYY-MM-DD HH:mm:ss. For example 'Home;2019-03-05 20:01:00;3'")
var ErrInvalidCapacity = errors.New("Capacity must be a number. Use 0 for unlimited. For example 'Home;2019-03-05 20:01:00;3'")
var defaultLocation *time.Location

const goingResponse = "OK, going!"
const notGoingResponse = "OK, not going :("
const goingCallbackData = "\fgoing"
const notGoingCallbackData = "\fnotGoing"
const goingIdentifier = "going"
const notGoingIdentifier = "notGoing"
const goingLabel = "Going"
const notGoingLabel = "Not going"
const nextEventTitle = "Next event!"
const nextEventDescription = "Where: %s, When: %s"
const meetingCreatedText = "Meeting created for %s at %s!"
const meetingCreatedDateFormat = "Monday 02 Jan 2006 15:04"
const invalidInputTitle = "Invalid input"

type editableMessage struct {
	MessageID string
	ChatID    int64
}

func (e *editableMessage) MessageSig() (messageID string, chatID int64) {
	messageID = e.MessageID
	chatID = e.ChatID
	return
}

func init() {
	var err error
	defaultLocation, err = time.LoadLocation("America/Argentina/Buenos_Aires") // FIXME: config timezone?
	if err != nil {
		panic(err)
	}
}

func meetingText(meeting *meetings.Meeting, users []*users.ExternalUser) string {
	usersText := ""
	if users != nil && len(users) > 0 {
		usersText = "\nAttendees:\n"
		for _, user := range users {
			usersText += fmt.Sprintf("* %s\n", user.DisplayName)
		}
	}
	return fmt.Sprintf(meetingCreatedText, meeting.Time.In(defaultLocation).Format(meetingCreatedDateFormat), meeting.Location) + usersText
}

func meetingOptions() *tb.SendOptions {
	goingButton := tb.InlineButton{
		Unique: goingIdentifier,
		Text:   goingLabel,
	}
	notGoingButton := tb.InlineButton{
		Unique: notGoingIdentifier,
		Text:   notGoingLabel,
	}
	return &tb.SendOptions{
		ReplyMarkup: &tb.ReplyMarkup{
			InlineKeyboard: [][]tb.InlineButton{
				[]tb.InlineButton{
					goingButton,
					notGoingButton,
				},
			},
		},
	}
}

func parseQuery(input string) (*meetings.Meeting, error) {
	data := strings.Split(input, ";")
	if len(data) != 3 {
		return nil, ErrNeedsSegments
	}
	date, err := time.Parse("2006-01-02 15:04:05", strings.TrimSpace(data[1]))
	if err != nil {
		return nil, ErrInvalidDate
	}
	capacity, err := strconv.Atoi(strings.TrimSpace(data[2]))
	if err != nil || capacity < 0 {
		return nil, ErrInvalidCapacity
	}
	return &meetings.Meeting{
		Time:     date,
		Location: strings.TrimSpace(data[0]),
		Capacity: capacity,
	}, nil
}

func formatUserDisplayName(user *tb.User) string {
	if user.Username != "" {
		if user.FirstName == "" && user.LastName == "" {
			return user.Username
		} else if user.FirstName == "" {
			return fmt.Sprintf("%s (%s)", user.Username, user.LastName)
		}
		return fmt.Sprintf("%s (%s %s)", user.Username, user.FirstName, user.LastName)
	}
	if user.FirstName == "" && user.LastName == "" {
		return "(unknown user)"
	} else if user.FirstName == "" {
		return user.LastName
	}
	return user.FirstName
}

func startTelegram(token string, mf *meetings.Factory, uf users.Factory) error {
	var b *tb.Bot
	b, err := tb.NewBot(tb.Settings{
		Token: token,
		Poller: tb.NewMiddlewarePoller(&tb.LongPoller{Timeout: 1 * time.Second}, func(upd *tb.Update) bool {
			if upd.Callback != nil {
				respond := func(text string) bool {
					err := b.Respond(upd.Callback, &tb.CallbackResponse{
						Text: text,
					})
					if err != nil {
						log.Print(err)
					}
					return false
				}
				respondEmpty := func() bool { return respond("") }

				if upd.Callback.Data != goingCallbackData && upd.Callback.Data != notGoingCallbackData {
					return true
				}

				going := upd.Callback.Data == goingCallbackData

				userID, err := uf.GetOrCreateUser(&users.ExternalUser{
					Source:      users.SourceTelegram,
					ID:          strconv.Itoa(upd.Callback.Sender.ID),
					DisplayName: formatUserDisplayName(upd.Callback.Sender),
				})
				if err != nil {
					return respondEmpty()
				}

				groupID, err := uf.GetOrCreateGroup(&users.ExternalGroup{Source: users.SourceTelegram, ID: strconv.FormatInt(upd.Callback.Message.Chat.ID, 10)})
				if err != nil {
					return respondEmpty()
				}

				if going {
					err = mf.UserRSVPMeeting(groupID, &meetings.Attendee{UserID: userID, Amount: 1})
				} else {
					err = mf.UserRSVPMeeting(groupID, &meetings.Attendee{UserID: userID, Amount: 1})
				}
				if err != nil {
					if err == meetings.NoActiveMeeting {
						return respond(err.Error())
					}
					return respondEmpty()
				}

				if going {
					respond(goingResponse)
				} else {
					respond(notGoingResponse)
				}

				meetingMessage := &editableMessage{}
				err = mf.GetMeetingAttendeesData(groupID, meetingMessage)
				if err != nil {
					log.Print(err)
					return false
				}
				if meetingMessage.MessageID != "" && meetingMessage.ChatID != 0 {
					meeting, err := mf.GetMeeting(groupID)
					if err != nil {
						log.Print(err)
						return false
					}
					attendees, err := mf.GetMeetingAttendees(groupID)
					if err != nil {
						log.Print(err)
						return false
					}
					attendeesUserID := make([]string, len(attendees))
					for i, att := range attendees {
						attendeesUserID[i] = att.UserID
					}
					usersMap, err := uf.GetUsers(attendeesUserID)
					if err != nil {
						log.Print(err)
						return false
					}
					users := make([]*users.ExternalUser, 0, len(usersMap))
					for _, user := range usersMap {
						users = append(users, user)
					}
					b.Edit(meetingMessage, meetingText(meeting, users), meetingOptions())
				}
				return false

			}
			return true
		}),
	})

	if err != nil {
		return err
	}

	b.Handle(tb.OnQuery, func(q *tb.Query) {
		m, err := parseQuery(q.Text)
		if err != nil {
			err = b.Answer(q, &tb.QueryResponse{
				Results: tb.Results{
					&tb.ArticleResult{
						Title:       invalidInputTitle,
						Description: err.Error(),
						Text:        "-",
					},
				},
				CacheTime: 60,
			})
			if err != nil {
				log.Print(err)
			}
			return
		}

		err = b.Answer(q, &tb.QueryResponse{
			Results: tb.Results{
				&tb.ArticleResult{
					Title:       nextEventTitle,
					Description: fmt.Sprintf(nextEventDescription, m.Location, m.Time.Format(time.RFC1123)),
					Text:        q.Text,
				},
			},
			CacheTime: 60,
		})

		if err != nil {
			log.Print(err)
		}

	})

	b.Handle(tb.OnText, func(m *tb.Message) {
		if !m.FromGroup() {
			return
		}
		meeting, err := parseQuery(m.Text)
		if err != nil {
			return
		}
		groupID, err := uf.GetOrCreateGroup(&users.ExternalGroup{
			Source: users.SourceTelegram,
			ID:     strconv.FormatInt(m.Chat.ID, 10),
		})
		if err != nil {
			return
		}
		err = mf.CreateMeeting(groupID, meeting)
		if err != nil {
			if err == meetings.MeetingAlreadyActive || err == meetings.MeetingIsInThePast {
				b.Send(m.Chat, err.Error())
			}
			return
		}
		message, err := b.Send(m.Chat, meetingText(meeting, nil), meetingOptions())
		if err != nil {
			log.Print(err)
			return
		}
		err = mf.SetMeetingAttendeesData(groupID, &editableMessage{
			MessageID: strconv.Itoa(message.ID),
			ChatID:    m.Chat.ID,
		})
		if err != nil {
			log.Print(err)
			return
		}
	})

	b.Start()
	return nil
}
