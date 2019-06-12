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

func init() {
	var err error
	defaultLocation, err = time.LoadLocation("America/Argentina/Buenos_Aires") // FIXME: config timezone?
	if err != nil {
		panic(err)
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

				if upd.Callback.Data != "\fgoing" && upd.Callback.Data != "\fnotGoing" {
					return true
				}

				going := upd.Callback.Data == "\fgoing"

				userID, err := uf.GetOrCreateUser(&users.ExternalUser{Source: users.SourceTelegram, ID: strconv.Itoa(upd.Callback.Sender.ID)})
				if err != nil {
					return respondEmpty()
				}

				groupID, err := uf.GetOrCreateGroup(&users.ExternalGroup{Source: users.SourceTelegram, ID: strconv.FormatInt(upd.Callback.Message.Chat.ID, 10)})
				if err != nil {
					return respondEmpty()
				}

				if going {
					err = mf.AddUserToMeeting(groupID, userID)
				} else {
					err = mf.RemoveUserFromMeeting(groupID, userID)
				}
				if err != nil {
					return respondEmpty()
				}

				if going {
					return respond("OK, going")
				}
				return respond("OK, not going")
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
						Title:       "Invalid input",
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
					Title:       "Next event!",
					Description: fmt.Sprintf("Where: %s, When: %s", m.Location, m.Time.Format(time.RFC1123)),
					Text:        q.Text,
				},
			},
			CacheTime: 60,
		})

		if err != nil {
			fmt.Println(err)
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
		goingButton := tb.InlineButton{
			Unique: "going",
			Text:   "Going",
		}
		notGoingButton := tb.InlineButton{
			Unique: "notGoing",
			Text:   "Not Going",
		}
		b.Send(m.Chat, fmt.Sprintf("Meeting created for %s at %s!", meeting.Time.In(defaultLocation).Format("Monday 02 Jan 2006 15:04"), meeting.Location), &tb.SendOptions{
			ReplyMarkup: &tb.ReplyMarkup{
				InlineKeyboard: [][]tb.InlineButton{
					[]tb.InlineButton{
						goingButton,
						notGoingButton,
					},
				},
			},
		})
	})

	b.Start()
	return nil
}
