// +build wasm,web

package fight

import (
	"bytes"
	"encoding/gob"
	"github.com/gontikr99/chutzparse/internal/model/fight"
	"github.com/gontikr99/chutzparse/pkg/console"
	"github.com/gontikr99/chutzparse/pkg/electron/ipc"
	"github.com/gontikr99/chutzparse/pkg/electron/ipc/ipcrenderer"
)

const maxRetention = 1000

var finishedFights []*fight.Fight
var listeners = map[int]chan struct{}{}
var listenerId = 0

func init() {
	fightChan, _ := ipc.GetChunkedEndpoint(ipcrenderer.Endpoint{}).Listen(fight.ChannelFinishedFights)
	go func() {
		for {
			fightIn := <-fightChan
			var fightData fight.Fight
			err := gob.NewDecoder(bytes.NewReader(fightIn.Content())).Decode(&fightData)
			if err == nil {
				finishedFights = append(finishedFights, &fightData)
				if len(finishedFights)>maxRetention {
					finishedFights = append([]*fight.Fight(nil), finishedFights[len(finishedFights)-maxRetention:]...)
				}
				for _, listener := range listeners {
					func() {
						defer func() { recover() }()
						listener <- struct{}{}
					}()
				}
			} else {
				console.Log(err)
			}
		}
	}()
}

func listenForFights() (<-chan struct{}, func()) {
	id := listenerId
	listenerId++
	signal := make(chan struct{})
	listeners[id] = signal
	doneFunc := func() {
		delete(listeners, id)
		close(signal)
	}
	return signal, doneFunc
}
