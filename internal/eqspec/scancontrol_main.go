// +build wasm,electron

package eqspec

import (
	"context"
	"net/rpc"
)

type scanCtlServer struct{}

func (s scanCtlServer) Restart() error {
	RestartLogScans(context.Background())
	return nil
}

func HandleRPC() func(server *rpc.Server) {
	return handleScanControl(scanCtlServer{})
}
