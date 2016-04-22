package tbt

import (
	tbot "github.com/go-telegram-bot-api/telegram-bot-api"
)

func GetUser(ctx *Context) *tbot.User {
	if ctx.upd.Message != nil {
		return ctx.upd.Message.From
	} else if ctx.upd.CallbackQuery != nil {
		return ctx.upd.CallbackQuery.Message.From
	} else if ctx.upd.ChosenInlineResult != nil {
		return ctx.upd.ChosenInlineResult.From
	} else if ctx.upd.InlineQuery != nil {
		return ctx.upd.ChosenInlineResult.From
	}
	return nil
}
