// +build wasm,electron

package parsecomms

import (
	"bytes"
	"encoding/gob"
	"github.com/gontikr99/chutzparse/internal/parse_model/parsedefs"
	"github.com/gontikr99/chutzparse/pkg/console"
	"github.com/gontikr99/chutzparse/pkg/electron/browserwindow"
)

// BroadcastHitEvent sends hit events to all windows.
func BroadcastHitEvent(channel string, msg *parsedefs.HitEvent) {
	buf := &bytes.Buffer{}
	err := gob.NewEncoder(buf).Encode(msg)
	if err == nil {
		browserwindow.Broadcast(channel, buf.Bytes())
	} else {
		console.Log(err)
	}
}

// BroadcastThroughput sends throughput bars to all windows
func BroadcastThroughput(displays []parsedefs.ThroughputState) {
	buf := &bytes.Buffer{}
	err := gob.NewEncoder(buf).Encode(&parsedefs.ThroughputStateEvent{Content: displays})
	if err == nil {
		browserwindow.Broadcast(parsedefs.ChannelThroughput, buf.Bytes())
	} else {
		console.Log(err)
	}
}