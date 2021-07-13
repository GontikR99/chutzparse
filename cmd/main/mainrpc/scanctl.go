// +build wasm,electron

package mainrpc

import (
	"context"
	"github.com/gontikr99/chutzparse/internal/eqspec"
)

type scanCtlServer struct {}

func (s scanCtlServer) Restart() error {
	eqspec.RestartLogScans(context.Background())
	return nil
}

func init() {
	register(eqspec.HandleScanControl(scanCtlServer{}))
}