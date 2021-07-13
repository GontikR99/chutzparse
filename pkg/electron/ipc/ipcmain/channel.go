// +build wasm,electron

package ipcmain

import (
	"github.com/gontikr99/chutzparse/pkg/electron"
	"github.com/gontikr99/chutzparse/pkg/electron/ipc"
	"github.com/gontikr99/chutzparse/pkg/jsbinding"
	"strconv"
	"syscall/js"
)

var ipcMain = electron.JSValue().Get("ipcMain")

func Listen(channelName string) (<-chan ipc.Message, func()) {
	resultChan := make(chan ipc.Message)
	recvFunc := js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		event := args[0]
		data := jsbinding.ReadArrayBuffer(args[1])
		resultChan <- &electronMessage{
			event:   event,
			content: []byte(data),
		}
		return nil
	})
	ipcMain.Call("on", ipc.Prefix+channelName, recvFunc)
	return resultChan, func() {
		ipcMain.Call("removeListener", ipc.Prefix+channelName, recvFunc)
		recvFunc.Release()
		close(resultChan)
	}
}

type electronMessage struct {
	event   js.Value
	content []byte
}

func (e *electronMessage) JSValue() js.Value {
	return e.event
}

func (e *electronMessage) Content() []byte {
	return e.content
}

func (e *electronMessage) Sender() string {
	return strconv.Itoa(e.event.Get("sender").Get("id").Int())
}

func (e *electronMessage) Reply(channelName string, data []byte) {
	e.event.Call("reply", ipc.Prefix+channelName, jsbinding.MakeArrayBuffer(data))
}
