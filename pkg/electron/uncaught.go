//go:build wasm && electron
// +build wasm,electron

package electron

import (
	"errors"
	"github.com/gontikr99/chutzparse/pkg/console"
	"github.com/gontikr99/chutzparse/pkg/nodejs"
	"syscall/js"
)

type CallbackHandle int

var handleGen = CallbackHandle(0)
var uncaughtCallbacks = map[CallbackHandle]func(value error){}
var electronAlert = nodejs.Require("electron-alert")

func init() {
	js.Global().Get("process").Call("on", "uncaughtException",
		js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
			if len(uncaughtCallbacks) != 0 {
				msg := args[0].Get("message")
				if args[0].IsUndefined() {
					msg = js.ValueOf("Unknown error")
				}
				console.LogRaw(args[0])
				for _, callback := range uncaughtCallbacks {
					callback(errors.New(msg.String()))
				}
			} else {
				electronAlert.Call("uncaughtException", false,
					js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
						console.LogRaw(args[0])
						return nil
					}),
					true, true).Invoke(args[0])
			}
			return nil
		}))
}

func RegisterUncaughtException(callback func(error)) CallbackHandle {
	h := handleGen
	handleGen++
	uncaughtCallbacks[h] = callback
	return h
}

func (h CallbackHandle) Release() {
	delete(uncaughtCallbacks, h)
}
