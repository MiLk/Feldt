package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/MiLk/feldt/currency"
	"github.com/MiLk/feldt/timer"
)

type Module interface {
	Start(*tgbotapi.BotAPI) error
	Stop()
	ProcessMessage(update tgbotapi.Update) (tgbotapi.Chattable, error)
	ProcessCallbackQuery(update tgbotapi.Update) (tgbotapi.Chattable, error)
}

var moduleList = []Module{
	currency.New(),
	timer.New(),
}

func StartAllModules(api *tgbotapi.BotAPI) error {
	for _, p := range moduleList {
		if err := p.Start(api); err != nil {
			return err
		}
	}

	return nil
}

func StopAllModules() {
	for _, p := range moduleList {
		p.Stop()
	}
}
