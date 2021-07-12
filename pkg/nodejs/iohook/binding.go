// +build wasm,electron

package iohook

import (
	"github.com/gontikr99/chutzparse/pkg/nodejs"
	"syscall/js"
)

var jsModule = nodejs.Require("iohook")

type KeyEvent struct {
	Type     string
	KeyCode  int
	RawCode  int
	AltKey   bool
	CtrlKey  bool
	ShiftKey bool
	MetaKey  bool
}

func on(eventType string, callback func(arg js.Value)) {
	jsModule.Call("on", eventType, js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		callback(args[0])
		return nil
	}))
}

func newKeyEvent(eventRaw js.Value) *KeyEvent {
	return &KeyEvent{
		Type:     eventRaw.Get("type").String(),
		KeyCode:  eventRaw.Get("keycode").Int(),
		RawCode:  eventRaw.Get("rawcode").Int(),
		AltKey:   eventRaw.Get("altKey").Bool(),
		CtrlKey:  eventRaw.Get("ctrlKey").Bool(),
		ShiftKey: eventRaw.Get("shiftKey").Bool(),
		MetaKey:  eventRaw.Get("metaKey").Bool(),
	}
}

func OnKeyDown(callback func(*KeyEvent)) {
	on("keydown", func(arg js.Value) { callback(newKeyEvent(arg)) })
}

func OnKeyUp(callback func(*KeyEvent)) {
	on("keyup", func(arg js.Value) { callback(newKeyEvent(arg)) })
}

func Start() {
	jsModule.Call("start", false)
}
