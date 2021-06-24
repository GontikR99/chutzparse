// +build wasm,electron

package eqwnd

import (
	"errors"
	"github.com/gontikr99/chutzparse/pkg/electron"
	"github.com/gontikr99/chutzparse/pkg/nodejs"
	"syscall/js"
)

var win32api=nodejs.Require("win32-api")
var user32=win32api.Get("U").Call("load")
var buffer=js.Global().Get("Buffer")

func findWindow() (uint32, error) {
	nameBuf := buffer.Call("from", "EverQuest\x00", "ucs2")
	resultInt := user32.Call("FindWindowExW", 0, 0, js.Null(), nameBuf).Int()
	if resultInt==0 {
		return 0, errors.New("window not found")
	} else {
		return uint32(resultInt), nil
	}
}

func IsTop() bool {
	hwnd, err := findWindow()
	if err!=nil {
		return false
	}
	topHwnd := uint32(user32.Call("GetForegroundWindow").Int())
	return hwnd==topHwnd
}

func GetExtents() (*electron.Rectangle, error) {
	hwnd, err := findWindow()
	if err!=nil {return nil, err}

	rectRef := buffer.Call("alloc", 16)
	success := user32.Call("GetWindowRect", hwnd, rectRef).Int()
	if success==0 {
		return nil, errors.New("failed to get window extents")
	}
	left := rectRef.Call("readInt32LE", 0).Int()
	top := rectRef.Call("readInt32LE", 4).Int()
	right := rectRef.Call("readInt32LE", 8).Int()
	bottom := rectRef.Call("readInt32LE", 12).Int()

	return &electron.Rectangle{
		X:      left,
		Y:      top,
		Width:  right-left,
		Height: bottom-top,
	}, nil
}