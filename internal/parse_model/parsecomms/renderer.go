// +build wasm,web

package parsecomms

import (
	"bytes"
	"context"
	"encoding/gob"
	"github.com/gontikr99/chutzparse/internal/parse_model/parsedefs"
	"github.com/gontikr99/chutzparse/pkg/console"
	"github.com/gontikr99/chutzparse/pkg/electron/ipc/ipcrenderer"
)

func HitDisplayListen(ctx context.Context, channel string) <-chan *parsedefs.HitEvent {
	chn, _ := ipcrenderer.Endpoint{}.Listen(channel)
	outChan := make(chan *parsedefs.HitEvent)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case inMsg := <-chn:
				hde := &parsedefs.HitEvent{}
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

func ThroughputListen(ctx context.Context) <-chan []parsedefs.ThroughputState {
	chn, _ := ipcrenderer.Endpoint{}.Listen(parsedefs.ChannelThroughput)
	outChan := make(chan []parsedefs.ThroughputState)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case inMsg := <-chn:
				tse := &parsedefs.ThroughputStateEvent{}
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