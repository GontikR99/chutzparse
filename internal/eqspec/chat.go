// +build wasm,electron

package eqspec

import (
	"github.com/gontikr99/chutzparse/pkg/multipattern"
)

type ChatLog struct {
	Source  string
	Text    string
	Method  ChatMethod
	Channel string
}

func (c *ChatLog) Visit(handler ParsedLogHandler) interface{} { return handler.OnChat(c) }

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
	MethodGlobal
)

func handleChat(mp *multipattern.Multipattern) *multipattern.Multipattern {
	return commonSubpatterns(mp).
		On("(.+) says, '(.*)'", func(parts []string, _ interface{}) interface{} {
			return &ChatLog{
				Source:  normalizeName(parts[1]),
				Text:    parts[2],
				Method:  MethodSay,
				Channel: "",
			}
		})
	// FIXME: add group, raid, guild, tell, shout, auction, ooc, global chat.  For both "you" and someone else.
}
