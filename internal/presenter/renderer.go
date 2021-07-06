// +build wasm,web

package presenter

import (
	"bytes"
	"context"
	"encoding/gob"
	"github.com/gontikr99/chutzparse/pkg/console"
	"github.com/gontikr99/chutzparse/pkg/electron/ipc/ipcrenderer"
)

func HitDisplayListen(ctx context.Context, channel string) <-chan *HitEvent {
	chn, _ := ipcrenderer.Endpoint{}.Listen(channel)
	outChan := make(chan *HitEvent)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case inMsg := <-chn:
				hde := &HitEvent{}
				err := gob.NewDecoder(bytes.NewReader(inMsg.Content())).Decode(hde)
				if err != nil {
					console.Log(err)
					continue
				}
				outChan <- hde
			}
		}
	}()
	return outChan
}

func ThroughputListen(ctx context.Context) <-chan []ThroughputState {
	chn, _ := ipcrenderer.Endpoint{}.Listen(ChannelThroughput)
	outChan := make(chan []ThroughputState)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case inMsg := <-chn:
				tse := &ThroughputStateEvent{}
				err := gob.NewDecoder(bytes.NewReader(inMsg.Content())).Decode(tse)
				if err != nil {
					console.Log(err)
					continue
				}
				outChan <- tse.Content
			}
		}
	}()
	return outChan
}
