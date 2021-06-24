// +build wasm,web

package electron

import "github.com/gontikr99/chutzparse/pkg/electron/ipc/ipcrenderer"

func IsPresent() bool {
	return ipcrenderer.Client != nil
}
