//go:build wasm && electron
// +build wasm,electron

package presenter

import (
	"bytes"
	"encoding/gob"
	"github.com/gontikr99/chutzparse/pkg/console"
	"github.com/gontikr99/chutzparse/pkg/electron/browserwindow"
)

// BroadcastHitEvent sends hit events to all windows.
func BroadcastHitEvent(channel string, msg *HitEvent) {
	buf := &bytes.Buffer{}
	err := gob.NewEncoder(buf).Encode(msg)
	if err == nil {
		browserwindow.Broadcast(channel, buf.Bytes())
	} else {
		console.Log(err)
	}
}

// BroadcastThroughput sends throughput bars to all windows
func BroadcastThroughput(displays []ThroughputState) {
	buf := &bytes.Buffer{}
	err := gob.NewEncoder(buf).Encode(&ThroughputStateEvent{Content: displays})
	if err == nil {
		browserwindow.BroadcastChunked(ChannelThroughput, buf.Bytes())
	} else {
		console.Log(err)
	}
}
