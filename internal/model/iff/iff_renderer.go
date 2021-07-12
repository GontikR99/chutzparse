// +build wasm,web

package iff

import (
	"bytes"
	"encoding/gob"
	"github.com/gontikr99/chutzparse/pkg/console"
	"github.com/gontikr99/chutzparse/pkg/electron/ipc/ipcrenderer"
)

// postUpdate sends a change made in a renderer to the main process.
func postUpdate(update IffUpdate) {
}


func init() {
	go func() {
		inChan, _ := ipcrenderer.Endpoint{}.Listen(channelIffUpdate)
		for {
			updateMsg := <-inChan
			var holder IffUpdateHolder
			err := gob.NewDecoder(bytes.NewReader(updateMsg.Content())).Decode(&holder)
			if err == nil {
				holder.Update.Apply()
			} else {
				console.Log(err)
			}
		}
	}()
}
