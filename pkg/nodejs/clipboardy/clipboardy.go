// +build wasm,electron

package clipboardy

import "github.com/gontikr99/chutzparse/pkg/nodejs"

var clipboardyJs = nodejs.Require("clipboardy")

func WriteSync(text string) {
	clipboardyJs.Call("writeSync", text)
}