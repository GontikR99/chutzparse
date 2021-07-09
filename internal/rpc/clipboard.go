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

type ClipboardClient struct {
	cs ClipboardServer
}

func (cc *ClipboardClient) Copy(req *ClipboardCopyRequest, res *ClipboardCopyResponse) error {
	return cc.cs.Copy(req.Text)
}

func CopyClipboard(client *rpc.Client, text string) error {
	req := &ClipboardCopyRequest{Text: text}
	res := new(ClipboardCopyResponse)
	return client.Call("ClipboardClient.Copy", req, res)
}

func HandleClipboard(clipboard ClipboardServer) func(server *rpc.Server) {
	cc := &ClipboardClient{clipboard}
	return func(server *rpc.Server) {
		server.Register(cc)
	}
}