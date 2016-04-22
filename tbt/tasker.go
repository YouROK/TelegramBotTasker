package tbt

import (
	"log"

	tbot "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Handler interface {
	//if return false then next step
	Handle(*Context) bool
}

type Tasker struct {
	bot      *tbot.BotAPI
	handlers []Handler
}

func NewTasker(tokenId string) *Tasker {
	t := &Tasker{}

	var err error

	t.bot, err = tbot.NewBotAPI(tokenId)
	if err != nil {
		log.Panic(err)
	}
	return t
}

func (m *Tasker) AddController(hndl Handler) {
	m.handlers = append(m.handlers, hndl)
}

func (m *Tasker) Start() {

	if len(m.handlers) == 0 {
		log.Panic("Error handlers is empty")
	}

	u := tbot.NewUpdate(0)
	u.Timeout = 60

	updates, err := m.bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		go m.update(&update)
	}
}

func (m *Tasker) update(update *tbot.Update) {
	ctx := NewContext(m.bot)
	ctx.upd = update
	for _, h := range m.handlers {
		if h.Handle(ctx) {
			return
		}
	}
	if ctx.msg != nil {
		m.bot.Send(ctx.msg)
	}
}
