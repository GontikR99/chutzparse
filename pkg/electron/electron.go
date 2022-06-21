//go:build wasm && electron
// +build wasm,electron

package electron

import (
	"github.com/gontikr99/chutzparse/pkg/nodejs"
	"syscall/js"
)

var electronJs = nodejs.Require("electron")

func JSValue() js.Value {
	return electronJs
}
