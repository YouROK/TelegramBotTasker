package tbt

import (
	"log"

	tbot "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Context struct {
	upd    *tbot.Update
	bot    *tbot.BotAPI
	param  map[string]interface{}
	msg    tbot.Chattable
	routes *Routes
}

func NewContext(bot *tbot.BotAPI) *Context {
	c := &Context{}
	c.bot = bot
	return c
}

func (c *Context) GetUpdate() *tbot.Update {
	return c.upd
}

func (c *Context) GetBotApi() *tbot.BotAPI {
	return c.bot
}

func (c *Context) NewMessage(txt string) *tbot.MessageConfig {
	msg := tbot.NewMessage(c.upd.Message.Chat.ID, txt)
	return &msg
}

func (c *Context) SendText(txt string) *tbot.Message {
	msg := tbot.NewMessage(c.upd.Message.Chat.ID, txt)
	msg.DisableNotification = true
	msg.DisableWebPagePreview = true
	msg.ParseMode = tbot.ModeHTML
	return c.Send(msg)
}

func (c *Context) Send(msg tbot.Chattable) *tbot.Message {
	m, err := c.bot.Send(msg)
	if err != nil {
		log.Println(err)
		return nil
	}
	return &m
}

func (c *Context) SetParam(name string, val interface{}) {
	if c.param == nil {
		c.param = make(map[string]interface{})
	}
	c.param[name] = val
}

func (c *Context) GetParam(name string) interface{} {
	if c.param == nil {
		return nil
	}
	return c.param[name]
}

func (c *Context) ResponseMessage() *tbot.MessageConfig {
	if _, ok := c.msg.(*tbot.MessageConfig); ok || c.msg == nil {
		msg := c.NewMessage("")
		msg.DisableNotification = true
		msg.DisableWebPagePreview = true
		msg.ParseMode = tbot.ModeHTML
		c.msg = msg
	}
	return c.msg.(*tbot.MessageConfig)
}

func (c *Context) SetResponse(msg tbot.Chattable) {
	c.msg = msg
}
