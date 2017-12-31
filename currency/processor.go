package currency

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func (p currencyModule) ProcessMessage(update tgbotapi.Update) (tgbotapi.Chattable, error) {
	m := currencyRegexp.FindStringSubmatch(update.Message.Text)
	if len(m) == 0 {
		return nil, nil
	}

	c := m[2]
	if !isLetter(c) {
		c = symbols[c]
	}
	c = strings.ToUpper(c)

	var buttons []tgbotapi.InlineKeyboardButton

	for _, currency := range currencies {
		buttons = append(buttons,
			tgbotapi.NewInlineKeyboardButtonData(
				currency,
				fmt.Sprintf("%s %s %s", m[1], c, currency),
			),
		)
	}

	msg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		fmt.Sprintf(
			"I can convert %s %s into a different currency. Choose one from the list.",
			m[1],
			c,
		),
	)
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(buttons)
	return &msg, nil
}

func (p currencyModule) ProcessCallbackQuery(update tgbotapi.Update) (tgbotapi.Chattable, error) {
	parts := strings.Split(update.CallbackQuery.Data, " ")

	baseValue, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return nil, err
	}

	msg := tgbotapi.NewMessage(
		update.CallbackQuery.Message.Chat.ID,
		fmt.Sprintf("%s %s -> %.2f %s",
			parts[0],
			parts[1],
			baseValue*p.poller.rates[parts[1]][parts[2]],
			parts[2],
		),
	)
	return &msg, nil
}
