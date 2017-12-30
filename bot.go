package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"strconv"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

var currencies = []string{"EUR", "JPY", "USD", "BRL"}
var symbols = map[string]string{
	"€":   "EUR",
	"¥":   "JPY",
	"円":   "JPY",
	"$":   "USD",
	"‎R$": "BRL",
}
var currencyRegexp = regexp.MustCompile(`(?i)(?:\s|^)(\d+)\s*(€|EUR|¥|円|JPY|\$|USD|‎R\$|BRL)(?:\s|,|\?|$)`)

func IsLetter(s string) bool {
	for _, r := range s {
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') {
			return false
		}
	}
	return true
}

func processMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

	m := currencyRegexp.FindStringSubmatch(update.Message.Text)
	if len(m) == 0 {
		return
	}

	c := m[2]
	fmt.Println(c)
	if !IsLetter(c) {
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
	bot.Send(msg)
}

func processCallbackQuery(bot *tgbotapi.BotAPI, update tgbotapi.Update, rates map[string]map[string]float64) {
	log.Printf("[%s] %s", update.CallbackQuery.From.UserName, update.CallbackQuery.Message.Text)
	log.Printf("[%s] %s", update.CallbackQuery.From.UserName, update.CallbackQuery.Data)

	parts := strings.Split(update.CallbackQuery.Data, " ")

	baseValue, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		log.Println(err)
		return
	}

	msg := tgbotapi.NewMessage(
		update.CallbackQuery.Message.Chat.ID,
		fmt.Sprintf("%s %s -> %.2f %s",
			parts[0],
			parts[1],
			baseValue*rates[parts[1]][parts[2]],
			parts[2],
		),
	)
	bot.Send(msg)
}
