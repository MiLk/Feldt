package currency

import (
	"regexp"
)

type currencyModule struct {
	poller poller
}

func New() *currencyModule {
	return &currencyModule{
		poller: newPoller(),
	}
}

func (p currencyModule) Start() error {
	p.poller.start()
	return nil
}

func (p currencyModule) Stop() {
	p.poller.stop()
}

var currencies = []string{"EUR", "JPY", "USD", "BRL"}
var symbols = map[string]string{
	"€":   "EUR",
	"¥":   "JPY",
	"円":   "JPY",
	"$":   "USD",
	"‎R$": "BRL",
}
var currencyRegexp = regexp.MustCompile(`(?i)(?:\s|^)(\d+)\s*(€|EUR|¥|円|JPY|\$|USD|‎R\$|BRL)(?:\s|,|\?|$)`)

func isLetter(s string) bool {
	for _, r := range s {
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') {
			return false
		}
	}
	return true
}
