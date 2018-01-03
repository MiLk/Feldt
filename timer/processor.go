package timer

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

var timerRegexp = regexp.MustCompile(`(?i)^((?:\d+[smhdw]{1})+)(?:\s+(.+))?$`)
var durationRegexp = regexp.MustCompile(`(?i)(\d+)([smhdw]{1})`)

var unitMap = map[string]time.Duration{
	"ns": time.Nanosecond,
	"us": time.Microsecond,
	"µs": time.Microsecond, // U+00B5 = micro symbol
	"μs": time.Microsecond, // U+03BC = Greek letter mu
	"ms": time.Millisecond,
	"s":  time.Second,
	"m":  time.Minute,
	"h":  time.Hour,
	"d":  time.Hour * 24,
	"w":  time.Hour * 24 * 7,
}

func (p timerModule) ProcessMessage(update tgbotapi.Update) (tgbotapi.Chattable, error) {
	if update.Message.Entities == nil || len(*update.Message.Entities) == 0 {
		return nil, nil
	}

	var command *tgbotapi.MessageEntity
	for _, e := range *update.Message.Entities {
		if e.Type != "bot_command" {
			continue
		}
		entityText := update.Message.Text[e.Offset : e.Offset+e.Length]
		if entityText == "/t" || entityText == "/timer" {
			command = &e
		}
	}
	if command == nil {
		return nil, nil
	}

	if len(update.Message.Text) <= command.Offset+command.Length+1 {
		return nil, errors.New("no parameters specified")
	}

	m := timerRegexp.FindStringSubmatch(update.Message.Text[command.Offset+command.Length+1:])
	if len(m) == 0 {
		return nil, errors.New("invalid parameters")
	}

	durations := durationRegexp.FindAllStringSubmatch(m[1], -1)

	duration := 0 * time.Second
	for _, d := range durations {
		di, err := strconv.Atoi(d[1])
		if err != nil {
			return nil, err
		}
		duration += unitMap[d[2]] * time.Duration(di)
	}

	t, err := p.storage.add(update.Message.Chat.ID, update.Message.From.UserName, duration, m[2])
	if err != nil {
		return nil, err
	}

	msg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		fmt.Sprintf(
			"[timer added] %s (%s - %s)",
			t.reason,
			t.duration,
			t.end.Format("2006-02-01 15:04:05 MST"),
		),
	)
	return &msg, nil
}

func (p timerModule) ProcessCallbackQuery(update tgbotapi.Update) (tgbotapi.Chattable, error) {
	return nil, nil
}
