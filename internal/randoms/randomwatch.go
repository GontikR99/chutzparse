//go:build wasm && electron
// +build wasm,electron

package randoms

import (
	"github.com/gontikr99/chutzparse/internal/eqspec"
	"github.com/gontikr99/chutzparse/pkg/electron/browserwindow"
	"time"
)

var currentRolls = []*randomEntry{}

type randomEntry struct {
	Timestamp time.Time
	Source    string
	Lbound    int32
	Ubound    int32
	Value     int32
}

func init() {
	eqspec.RegisterLogsListener(func(entries []*eqspec.LogEntry) {
		didWork := false
		for _, entry := range entries {
			if msg, ok := entry.Meaning.(*eqspec.RandomLog); ok {
				currentRolls = append(currentRolls, &randomEntry{
					Timestamp: time.Now(),
					Source:    msg.Source,
					Lbound:    msg.Lbound,
					Ubound:    msg.Ubound,
					Value:     msg.Value,
				})
				didWork = true
			}
		}
		if didWork {
			browserwindow.Broadcast(ChannelChange, []byte{})
		}
	})
}
