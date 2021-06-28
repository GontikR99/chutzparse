// +build wasm,electron

package dmodel

import (
	"bytes"
	"encoding/gob"
	"github.com/gontikr99/chutzparse/internal/eqlog"
	"github.com/gontikr99/chutzparse/pkg/console"
	"github.com/gontikr99/chutzparse/pkg/electron/browserwindow"
	"strconv"
)

type hitDisplayProduceReceiver struct {}

func broadcastHitDisplay(channel string, msg *HitDisplayEvent) {
	buf := &bytes.Buffer{}
	err := gob.NewEncoder(buf).Encode(msg)
	if err==nil {
		browserwindow.Broadcast(channel, buf.Bytes())
	} else {
		console.Log(err)
	}
}

func renderAmt(amount int64) string {
	if amount<1000 {
		return strconv.FormatInt(amount, 10)
	} else if amount<1000000 {
		amtFlt := float64(amount)/1000.0
		return strconv.FormatFloat(amtFlt, 'g', 3, 64)+"k"
	} else {
		amtFlt := float64(amount)/1000000.0
		return strconv.FormatFloat(amtFlt, 'g', 3, 64)+"M"
	}
}

func ListenForHits() {
	eqlog.RegisterLogsListener(func(entries []*eqlog.LogEntry) {
		for _, entry := range entries {
			if dmgEntry, ok := entry.Meaning.(*eqlog.DamageLog); ok {
				if dmgEntry.Source == entry.Character {
					if dmgEntry.SpellName!="" {
						broadcastHitDisplay(ChannelTopTarget, &HitDisplayEvent{
							Text:  renderAmt(dmgEntry.Amount),
							Color: "yellow",
						})
					} else if dmgEntry.Element==eqlog.PhysicalDamage {
						broadcastHitDisplay(ChannelTopTarget, &HitDisplayEvent{
							Text:  renderAmt(dmgEntry.Amount),
							Color: "white",
						})
					} else {
						broadcastHitDisplay(ChannelTopTarget, &HitDisplayEvent{
							Text:  renderAmt(dmgEntry.Amount),
							Color: "gray",
						})
					}
				}
				if dmgEntry.Target == entry.Character {
					broadcastHitDisplay(ChannelBottomTarget, &HitDisplayEvent{
						Text:  renderAmt(dmgEntry.Amount),
						Color: "red",
					})
				}
			} else if healEntry, ok := entry.Meaning.(*eqlog.HealLog); ok {
				if healEntry.Source == entry.Character {
					broadcastHitDisplay(ChannelTopTarget, &HitDisplayEvent{
						Text:  renderAmt(healEntry.Actual),
						Color: "green",
					})
				}
				if healEntry.Target == entry.Character {
					broadcastHitDisplay(ChannelBottomTarget, &HitDisplayEvent{
						Text:  renderAmt(healEntry.Actual),
						Color: "green",
					})
				}
			}
		}
	})
}