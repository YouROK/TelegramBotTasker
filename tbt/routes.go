package tbt

import (
	"reflect"
	"strings"

	tbot "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Route interface {
	Start(*Context) bool
	Handler
}

type Routes struct {
	cmd          string
	defRouteName string
	useKeyb      bool

	usersRoute map[int]Route
	routes     map[string]Route
}

func NewRoutes(cmd string) *Routes {
	r := &Routes{}
	r.usersRoute = make(map[int]Route)
	r.routes = make(map[string]Route)
	r.cmd = cmd
	return r
}

func (r *Routes) AddRoute(name string, route Route) {
	r.routes[strings.ToLower(name)] = route
}

func (r *Routes) SetDefRoute(name string) {
	if r.findRoute(name) == nil {
		panic("Error default router not found")
	}
	r.defRouteName = name
}

func (r *Routes) UseKeyboard(use bool) {
	r.useKeyb = use
}

func (r *Routes) Handle(ctx *Context) bool {
	//Change routes for user, cmd set route
	if ctx.upd.Message != nil && ctx.upd.Message.Command() == r.cmd && r.cmd != "" {
		if rn := ctx.upd.Message.CommandArguments(); rn == "" {
			r.help(ctx)
		} else {
			rn = strings.ToLower(rn)
			if route := r.findRoute(rn); route == nil {
				r.help(ctx)
			} else {
				usr := GetUser(ctx)
				r.usersRoute[usr.ID] = route
				return route.Start(ctx)
			}
		}
	} else {
		//Route user
		rname := r.GetCurrentRoute(ctx)
		if rname != "" {
			if usrRoute, ok := r.routes[rname]; ok {
				return invokeHandle(ctx, usrRoute)
			}
		}
		r.help(ctx)
	}
	return true
}

func (r *Routes) findRoute(name string) Route {
	name = strings.ToLower(name)
	if route, ok := r.routes[name]; ok {
		return route
	}
	return nil
}

func (r *Routes) help(ctx *Context) {
	var keys [][]tbot.KeyboardButton
	txt := ""
	for k, _ := range r.routes {
		if r.useKeyb {
			keys = append(keys, tbot.NewKeyboardButtonRow(tbot.NewKeyboardButton("/"+r.cmd+" "+k)))
		}
		txt += "/" + r.cmd + " " + k + "\n"
	}

	msg := ctx.NewMessage(txt)
	msg.DisableWebPagePreview = true
	msg.DisableNotification = true
	msg.ParseMode = tbot.ModeHTML
	if r.useKeyb && !r.onceRoute() {
		rplkey := tbot.NewReplyKeyboard(keys...)
		rplkey.OneTimeKeyboard = true
		msg.ReplyMarkup = rplkey
	}
	ctx.Send(msg)
}

func (r *Routes) onceRoute() bool {
	return r.cmd == "" || len(r.routes) == 1
}

func (r *Routes) GetCurrentRoute(ctx *Context) string {
	if r.onceRoute() {
		for k, _ := range r.routes {
			return k
		}
	}

	usr := GetUser(ctx)
	if usrRoute, ok := r.usersRoute[usr.ID]; ok {
		for k, v := range r.routes {
			if v == usrRoute {
				return k
			}
		}
		return ""
	}
	return r.defRouteName
}

func (r *Routes) GetUsersRoute() map[int]Route {
	return r.usersRoute
}

func (r *Routes) GetRoutes() map[string]Route {
	return r.routes
}

func invokeHandle(ctx *Context, route Route) bool {
	upd := ctx.GetUpdate()

	//Command handlers
	if upd.Message != nil && upd.Message.Command() != "" {
		type temper interface {
			OnCommand(*Context, string, string) bool
		}
		if hndl, ok := route.(temper); ok {
			//general function exist
			return hndl.OnCommand(ctx, upd.Message.Command(), upd.Message.CommandArguments())
		} else {
			//find on__cmds__ in handler
			cmd := upd.Message.Command()
			if mth := findCmdMethod(ctx, route, cmd); mth != nil {
				return mth.Func.Call([]reflect.Value{reflect.ValueOf(route), reflect.ValueOf(ctx), reflect.ValueOf(upd.Message.Command()), reflect.ValueOf(upd.Message.CommandArguments())})[0].Bool()
			}
		}
	}

	if upd.CallbackQuery != nil {
		type temper interface {
			OnCallbackQuery(*Context, *tbot.CallbackQuery) bool
		}
		if hndl, ok := route.(temper); ok {
			return hndl.OnCallbackQuery(ctx, upd.CallbackQuery)
		}
	}

	if upd.ChosenInlineResult != nil {
		type temper interface {
			OnChosenInlineResult(*Context, *tbot.ChosenInlineResult) bool
		}
		if hndl, ok := route.(temper); ok {
			return hndl.OnChosenInlineResult(ctx, upd.ChosenInlineResult)
		}
	}

	if upd.InlineQuery != nil {
		type temper interface {
			OnInlineQuery(*Context, *tbot.InlineQuery) bool
		}
		if hndl, ok := route.(temper); ok {
			return hndl.OnInlineQuery(ctx, upd.InlineQuery)
		}
	}

	if upd.Message != nil {
		type temper interface {
			OnMessage(*Context, *tbot.Message) bool
		}
		if hndl, ok := route.(temper); ok {
			return hndl.OnMessage(ctx, upd.Message)
		}
	}
	return route.Handle(ctx)
}

// OnHelp(ctx *Context, cmd, arg string) bool
func findCmdMethod(ctx *Context, str interface{}, methodName string) *reflect.Method {
	mth := getMethodByName(str, methodName)
	if mth != nil && mth.Type.NumIn() == 4 {
		if (mth.Type.In(1).Kind() == reflect.Ptr && mth.Type.In(1) == reflect.ValueOf(ctx).Type()) &&
			(mth.Type.In(2).Kind() == reflect.String && mth.Type.In(3).Kind() == reflect.String) {
			return mth
		}
	}
	return nil
}

func getMethodByName(str interface{}, methodName string) *reflect.Method {
	stType := reflect.TypeOf(str)
	methodName = "on" + strings.ToLower(methodName)
	for i := 0; i < stType.NumMethod(); i++ {
		if mth := stType.Method(i); strings.HasPrefix(mth.Name, "On") && strings.ToLower(mth.Name) == methodName {
			if mth.Type.NumOut() == 1 && mth.Type.Out(0).Kind() == reflect.Bool {
				return &mth
			}
		}
	}
	return nil
}
