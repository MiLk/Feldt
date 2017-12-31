package timer

import (
	"fmt"
	"sort"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type timer struct {
	chatID   int64
	start    time.Time
	end      time.Time
	duration time.Duration
	reason   string
}

type storage struct {
	keys []int
	data map[int]timer

	stopCh   chan bool
	doneCh   chan bool
	updateCh chan bool
	addCh    chan timer
}

func newStorage() storage {
	return storage{
		data:     map[int]timer{},
		keys:     []int{},
		stopCh:   make(chan bool, 1),
		doneCh:   make(chan bool, 1),
		updateCh: make(chan bool, 1),
		addCh:    make(chan timer),
	}
}

func (s storage) add(chatID int64, d time.Duration, r string) (*timer, error) {
	n := time.Now()
	t := timer{
		chatID:   chatID,
		start:    n,
		end:      n.Add(d),
		duration: d,
		reason:   r,
	}
	s.addCh <- t
	return &t, nil
}

func (s storage) start(bot *tgbotapi.BotAPI) {
	defer close(s.doneCh)
	go func(b *tgbotapi.BotAPI) {
		timer := time.NewTimer(1 * time.Hour)
		timer.Stop()
		isActive := false

		for {
			select {
			case t := <-s.addCh:
				k := int(t.end.Unix())
				s.keys = append(s.keys, k)
				sort.Ints(s.keys)
				s.data[k] = t
				s.updateCh <- true
			case <-s.updateCh:
				if len(s.keys) == 0 {
					continue
				}
				if isActive && !timer.Stop() {
					<-timer.C
				}
				nextEvent := s.data[s.keys[0]]
				d := nextEvent.end.Sub(time.Now())
				timer.Reset(d)
				isActive = true
			case <-timer.C:
				isActive = false
				n := time.Now()
				toDelete := 0
				for _, k := range s.keys {
					evt := s.data[k]
					if evt.end.After(n) {
						break
					}
					b.Send(tgbotapi.NewMessage(
						evt.chatID,
						fmt.Sprintf("ding! %s", evt.reason),
					))
					delete(s.data, k)
					toDelete++
				}
				if toDelete > 0 {
					if toDelete >= len(s.keys) {
						s.keys = []int{}
					} else {
						s.keys = s.keys[toDelete-1:]
						sort.Ints(s.keys)
					}
				}
				s.updateCh <- true
			case <-s.stopCh:
				fmt.Println("stop")
				return
			}
		}
	}(bot)
}

func (s storage) stop() {
	close(s.stopCh)
	close(s.updateCh)
	<-s.doneCh
}
