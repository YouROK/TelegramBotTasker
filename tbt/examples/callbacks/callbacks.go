package main

import (
	"log"
	"strconv"
	"tbt"

	tbot "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	routes *tbt.Routes
)

type Echo struct {
	count int
}

func (*Echo) Start(ctx *tbt.Context) bool {
	ctx.SendText("This is echo")
	return true
}

func (*Echo) OnMessage(ctx *tbt.Context, msg *tbot.Message) bool {
	resp := ctx.ResponseMessage()
	resp.Text = msg.Text
	buttons := tbot.NewInlineKeyboardRow(tbot.NewInlineKeyboardButtonData("Echo", msg.Text))
	keyb := tbot.NewInlineKeyboardMarkup(buttons)
	resp.ReplyMarkup = keyb
	return false
}

func (e *Echo) OnCallbackQuery(ctx *tbt.Context, cbq *tbot.CallbackQuery) bool {
	e.count++
	resp := tbot.NewEditMessageText(cbq.Message.Chat.ID, cbq.Message.MessageID, "Edit: "+strconv.Itoa(e.count))
	buttons := tbot.NewInlineKeyboardRow(tbot.NewInlineKeyboardButtonData("Echo "+strconv.Itoa(e.count), cbq.Data))
	keyb := tbot.NewInlineKeyboardMarkup(buttons)
	resp.ReplyMarkup = &keyb
	ctx.SetResponse(resp)
	return false
}

func (*Echo) Handle(ctx *tbt.Context) bool {
	log.Println(ctx.GetUpdate().ChosenInlineResult)
	log.Println(ctx.GetUpdate().InlineQuery)
	log.Println(ctx.GetUpdate().CallbackQuery)
	log.Println(ctx.GetUpdate().Message)
	return false
}

//////////////////////////////////////////////////////////

func main() {
	log.Println("Start")
	t := tbt.NewTasker("Your token")
	routes = tbt.NewRoutes("")
	routes.AddRoute("echo", &Echo{})
	routes.UseKeyboard(true)

	t.AddController(routes)
	t.Start()
}
