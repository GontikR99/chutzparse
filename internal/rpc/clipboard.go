package rpc

import (
	"net/rpc"
)

type ClipboardCopyRequest struct {
	Text string
}

type ClipboardCopyResponse struct {
}

type ClipboardServer interface {
	Copy(text string) error
}

type StubClipboard struct {
	cs ClipboardServer
}

func (cc *StubClipboard) Copy(req *ClipboardCopyRequest, res *ClipboardCopyResponse) error {
	return cc.cs.Copy(req.Text)
}

func CopyClipboard(client *rpc.Client, text string) error {
	req := &ClipboardCopyRequest{Text: text}
	res := new(ClipboardCopyResponse)
	return client.Call("StubClipboard.Copy", req, res)
}

func HandleClipboard(clipboard ClipboardServer) func(server *rpc.Server) {
	cc := &StubClipboard{clipboard}
	return func(server *rpc.Server) {
		server.Register(cc)
	}
}