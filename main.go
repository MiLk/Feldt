package main

import (
	"log"

	"os"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	token := os.Getenv("TELEGRAM_API_TOKEN")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	p := newPoller()
	p.start()
	defer p.stop()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	// Optional: wait for updates and clear them if you don't want to handle
	// a large backlog of old messages
	//time.Sleep(time.Millisecond * 500)
	//updates.Clear()

	for update := range updates {
		if update.CallbackQuery != nil {
			processCallbackQuery(bot, update, p.rates)
			continue
		}

		if update.Message != nil {
			processMessage(bot, update)
		}

	}
}
