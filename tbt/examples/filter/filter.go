package main

import (
	"log"
	"tbt"
)

var (
	routes *tbt.Routes
)

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

//////////////////////////////////////////////////////////

func main() {
	log.Println("Start")
	t := tbt.NewTasker("Your token")

	t.AddController(tbt.NewFilterMessage(1, 3))

	routes = tbt.NewRoutes("")
	routes.AddRoute("echo", &Echo{})
	routes.UseKeyboard(true)
	t.AddController(routes)
	t.Start()
}
