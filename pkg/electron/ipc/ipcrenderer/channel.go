// +build wasm

package ipcrenderer

import (
	"github.com/gontikr99/chutzparse/pkg/electron/ipc"
	"github.com/gontikr99/chutzparse/pkg/jsbinding"
	"net/rpc"
	"syscall/js"
)

var ipcRenderer = js.Global().Get("ipcRenderer")

type Endpoint struct{}

func (i Endpoint) Listen(channelName string) (<-chan ipc.Message, func()) {
	resultChan := make(chan ipc.Message, 16)
	recvFunc := js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		event := args[0]
		data := jsbinding.ReadArrayBuffer(args[1])
		func() {
			defer func() { recover() }()
			resultChan <- &electronMessage{
				event:   event,
				content: data,
			}
		}()
		return nil
	})
	ipcRenderer.Call("on", ipc.Prefix+channelName, recvFunc)
	return resultChan, func() {
		ipcRenderer.Call("removeListener", ipc.Prefix+channelName, recvFunc)
		recvFunc.Release()
		close(resultChan)
	}
}

func (i Endpoint) Send(channelName string, content []byte) {
	ipcRenderer.Call("send", ipc.Prefix+channelName, jsbinding.MakeArrayBuffer(content))
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
	return "mainProcess"
}

func (e *electronMessage) Reply(channelName string, data []byte) {
	e.event.Call("reply", ipc.Prefix+channelName, jsbinding.MakeArrayBuffer(data))
}

// Renderer side RPC client
var Client *rpc.Client

func init() {
	if !ipcRenderer.IsUndefined() {
		Client = ipc.NewClient(ipc.ChannelRPCMain, &Endpoint{})
	}
}
