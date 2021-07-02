// +build wasm,electron

package eqlog

import (
	"github.com/gontikr99/chutzparse/pkg/multipattern"
)

type ChatLog struct {
	Source string
	Text string
	Method ChatMethod
	Channel string
}

func (c *ChatLog) Visit(handler ParsedLogHandler) interface{} {return handler.OnChat(c)}

type ChatMethod int
const (
	MethodSay = ChatMethod(iota)
	MethodGroup
	MethodRaid
	MethodGuild
	MethodTell
	MethodShout
	MethodAuction
	MethodOOC
)

func handleChat(mp *multipattern.Multipattern) *multipattern.Multipattern {
	return commonSubpatterns(mp).
		On("(.+) says, '(.*)'", func(parts []string) interface{} {
			return &ChatLog{
				Source:  normalizeName(parts[1]),
				Text:    parts[2],
				Method:  MethodSay,
				Channel: "",
			}
		})
}