package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/MiLk/feldt/currency"
)

type Module interface {
	Start() error
	Stop()
	ProcessMessage(update tgbotapi.Update) (tgbotapi.Chattable, error)
	ProcessCallbackQuery(update tgbotapi.Update) (tgbotapi.Chattable, error)
}

var moduleList = []Module{
	currency.New(),
}

func StartAllModules() error {
	for _, p := range moduleList {
		if err := p.Start(); err != nil {
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
