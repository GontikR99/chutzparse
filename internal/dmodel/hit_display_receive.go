// +build wasm,web

package dmodel

import (
	"bytes"
	"encoding/gob"
	"github.com/gontikr99/chutzparse/pkg/electron/ipc/ipcrenderer"
)

func HitDisplayListen(channel string, callback func(event *HitDisplayEvent)) {
	chn, _ := ipcrenderer.Endpoint{}.Listen(channel)
	go func() {
		for {
			inMsg := <-chn
			hde := &HitDisplayEvent{}
			err := gob.NewDecoder(bytes.NewReader(inMsg.Content())).Decode(hde)
			if err==nil {
				callback(hde)
			}
		}
	}()
}