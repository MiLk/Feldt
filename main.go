package main

import (
	"log"
	"os"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	if err := mainE(); err != nil {
		log.Panic(err)
	}
}

func mainE() error {
	token := os.Getenv("TELEGRAM_API_TOKEN")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return err
	}

	if os.Getenv("DEBUG") != "" {
		bot.Debug = true
	}

	log.Printf("Authorized on account @%s (%s)", bot.Self.UserName, bot.Self.FirstName)

	if err := StartAllModules(); err != nil {
		return err
	}
	defer StopAllModules()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		return err
	}

	// Optional: wait for updates and clear them if you don't want to handle
	// a large backlog of old messages
	//time.Sleep(time.Millisecond * 500)
	//updates.Clear()

UpdateLoop:
	for update := range updates {
		if update.CallbackQuery != nil &&
			update.CallbackQuery.Message.From.ID == bot.Self.ID {
			log.Printf("[%s] callback query: %s - %s",
				update.CallbackQuery.From.UserName,
				update.CallbackQuery.Message.Text,
				update.CallbackQuery.Data)

			for _, p := range moduleList {
				if m, err := p.ProcessCallbackQuery(update); err != nil {
					log.Println(err)
					continue UpdateLoop
				} else if m != nil {
					bot.Send(m)
					continue UpdateLoop
				}
			}
		}

		if update.Message != nil {
			log.Printf("[%s] message: %s",
				update.Message.From.UserName,
				update.Message.Text)

			for _, p := range moduleList {
				if m, err := p.ProcessMessage(update); err != nil {
					log.Println(err)
					continue UpdateLoop
				} else if m != nil {
					bot.Send(m)
					continue UpdateLoop
				}
			}
		}
	}

	return nil
}
