// +build wasm, electron

package mainrpc

import (
	"github.com/gontikr99/chutzparse/internal/rpc"
	"github.com/gontikr99/chutzparse/pkg/nodejs/clipboardy"
)

type clipboardServer struct {}

func (c clipboardServer) Copy(text string) error {
	clipboardy.WriteSync(text)
	return nil
}

func init() {
	register(rpc.HandleClipboard(clipboardServer{}))
}