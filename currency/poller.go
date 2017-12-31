package currency

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type poller struct {
	rates  map[string]map[string]float64
	stopCh chan bool
	doneCh chan bool
}

func newPoller() poller {
	return poller{
		rates:  map[string]map[string]float64{},
		stopCh: make(chan bool, 1),
		doneCh: make(chan bool, 1),
	}
}

func (p poller) start() {
	defer close(p.doneCh)
	go func() {
		for {
			for _, c := range currencies {
				if r, err := getExchangeRates(c); err != nil {
					log.Println(err)
				} else {
					p.rates[c] = r
				}
			}

			select {
			case <-time.After(1 * time.Hour):
				continue
			case <-p.stopCh:
				return
			}
		}
	}()
}

func (p poller) stop() {
	close(p.stopCh)
	<-p.doneCh
}

type ratesReponse struct {
	Base  string
	Date  string
	Rates map[string]float64
}

func getExchangeRates(base string) (map[string]float64, error) {
	r, err := http.Get("https://api.fixer.io/latest?base=" + base + "&symbols=" + strings.Join(currencies, ","))
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	var body ratesReponse
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}
	fmt.Println(body)

	return body.Rates, nil
}
