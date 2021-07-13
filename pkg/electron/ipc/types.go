// +build wasm

package ipc

import "syscall/js"

const Prefix = "golang-msgcomm-"

// A message passed between processes
type Message interface {
	// Content of the message
	Content() []byte

	// Unique identifer of the message's sender
	Sender() string

	// Send a response back to the sender
	Reply(channelName string, data []byte)

	JSValue() js.Value
}

type Sender interface {
	Send(channelName string, data []byte)
}

type Listener interface {
	Listen(channelName string) (recv <-chan Message, done func())
}

type Endpoint interface {
	Sender
	Listener
}
