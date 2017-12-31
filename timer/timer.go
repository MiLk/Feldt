package timer

import "github.com/go-telegram-bot-api/telegram-bot-api"

type timerModule struct {
	storage storage
}

func New() *timerModule {
	return &timerModule{
		storage: newStorage(),
	}
}

func (p timerModule) Start(bot *tgbotapi.BotAPI) error {
	p.storage.start(bot)
	return nil
}

func (p timerModule) Stop() {
	p.storage.stop()
}
