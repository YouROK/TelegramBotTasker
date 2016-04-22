package tbt

import (
	"sync"
	"time"
)

type filter struct {
	lastTime time.Time
	lastMsg  string
}

type filterMessage struct {
	timefilter time.Duration
	timeclean  time.Duration
	messages   map[int]*filter
	rwmut      sync.RWMutex

	isclean  bool
	cleanmut sync.Mutex
}

/*
timeFilter - time during which repeated messages will be filtered
timeClean - time after which will delete older messages
*/
func NewFilterMessage(timeFilter, timeClean time.Duration) *filterMessage {
	f := &filterMessage{}
	f.timefilter = timeFilter
	f.timeclean = timeClean
	f.isclean = false
	f.messages = make(map[int]*filter)
	return f
}

func (f *filterMessage) Handle(ctx *Context) bool {
	if f.timeclean == 0 {
		return true
	}
	f.rwmut.RLock()
	defer f.rwmut.RUnlock()
	msg := ctx.GetUpdate().Message
	if msg != nil {
		defer func() { go f.cleanup() }()
		id := ctx.GetUpdate().Message.From.ID
		lastmsg, ok := f.messages[id]
		if ok && lastmsg.lastMsg == msg.Text && lastmsg.lastTime.Add(time.Second*f.timefilter).After(time.Now()) {
			f.messages[id].lastTime = time.Now()
			return true
		} else {
			f.messages[id] = &filter{lastMsg: msg.Text, lastTime: time.Now()}
		}
	}
	return false
}

func (f *filterMessage) cleanup() {
	if !f.isclean && f.timeclean > 0 {
		f.isclean = true
		f.cleanmut.Lock()
		if f.isclean {
			for k, v := range f.messages {
				if v.lastTime.Add(time.Second * f.timeclean).Before(time.Now()) {
					f.rwmut.Lock()
					delete(f.messages, k)
					f.rwmut.Unlock()
				}
			}
		}
		f.isclean = false
		f.cleanmut.Unlock()
	}
}
