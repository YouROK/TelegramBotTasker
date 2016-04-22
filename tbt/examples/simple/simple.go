package main

import (
	"log"
	"tbt"

	tbot "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	routes *tbt.Routes
)

type Logger struct {
}

func (tp *Logger) Handle(ctx *tbt.Context) bool {
	if ctx.GetUpdate().Message != nil {
		log.Println(ctx.GetUpdate().Message.From)
		log.Println(ctx.GetUpdate().Message.Text)
	}
	return false
}

type Echo struct {
}

func (*Echo) Start(ctx *tbt.Context) bool {
	ctx.SendText("This is echo")
	return true
}

func (*Echo) Handle(ctx *tbt.Context) bool {
	msg := ctx.GetUpdate().Message
	if msg != nil {
		ctx.SendText(msg.Text)
	}
	return true
}

type BlaBla struct {
}

func (*BlaBla) Start(ctx *tbt.Context) bool {
	resp := ctx.Response()
	resp.Text = "Bla bla bla"
	return false
}

func (*BlaBla) Handle(ctx *tbt.Context) bool {
	msg := ctx.GetUpdate().Message
	if msg != nil {
		resp := ctx.Response()
		resp.Text = msg.Text + "... bla bla bla"
		resp.ReplyToMessageID = msg.MessageID
	}
	return false
}

type MenuBla struct {
}

func (*MenuBla) Handle(ctx *tbt.Context) bool {
	if routes.GetCurrentRoute(ctx) == "bla" {
		msg := ctx.Response()
		msg.ReplyMarkup = tbot.NewReplyKeyboard(tbot.NewKeyboardButtonRow(tbot.NewKeyboardButton("Get all user tokens")),
			tbot.NewKeyboardButtonRow(tbot.NewKeyboardButton("Send email to user")))
	}
	return false
}

//////////////////////////////////////////////////////////

func main() {
	log.Println("Start")
	t := tbt.NewTasker("Your token")
	//Add logger
	t.AddController(&Logger{})
	//Setup routes, cmd ("set") never not send to route
	routes = tbt.NewRoutes("set")
	routes.AddRoute("echo", &Echo{})
	routes.AddRoute("bla", &BlaBla{})
	//	routes.SetDefRoute("echo")
	routes.UseKeyboard(true)

	t.AddController(routes)
	t.AddController(&MenuBla{})

	//if controller return true other controllers not execute
	t.Start()
}
