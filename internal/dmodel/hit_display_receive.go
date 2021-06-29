// +build wasm,web

package dmodel

import (
	"bytes"
	"encoding/gob"
	"github.com/gontikr99/chutzparse/pkg/console"
	"github.com/gontikr99/chutzparse/pkg/electron/ipc/ipcrenderer"
)

func HitDisplayListen(channel string) <-chan *HitDisplayEvent {
	chn, _ := ipcrenderer.Endpoint{}.Listen(channel)
	outChan := make(chan *HitDisplayEvent)
	go func() {
		for {
			inMsg := <-chn
			hde := &HitDisplayEvent{}
			err := gob.NewDecoder(bytes.NewReader(inMsg.Content())).Decode(hde)
			if err!=nil {
				console.Log(err)
				continue
			}
			outChan <- hde
		}
	}()
	return outChan
}