// +build wasm

package ipc

import (
	"github.com/gontikr99/chutzparse/pkg/jsbinding"
	"syscall/js"
)

func Encode(data []byte) js.Value {
	return jsbinding.MakeArrayBuffer(data)
}

func Decode(encoded js.Value) ([]byte,error) {
	return jsbinding.ReadArrayBuffer(encoded), nil
}