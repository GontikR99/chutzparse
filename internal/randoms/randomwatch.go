//go:build wasm && electron
// +build wasm,electron

package randoms

import (
	"github.com/gontikr99/chutzparse/internal/eqspec"
	"github.com/gontikr99/chutzparse/pkg/electron/browserwindow"
)

var currentRolls = []*eqspec.RandomLog{}

func init() {
	eqspec.RegisterLogsListener(func(entries []*eqspec.LogEntry) {
		didWork := false
		for _, entry := range entries {
			if msg, ok := entry.Meaning.(*eqspec.RandomLog); ok {
				currentRolls = append(currentRolls, msg)
				didWork = true
			}
		}
		if didWork {
			browserwindow.Broadcast(ChannelChange, []byte{})
		}
	})
}
