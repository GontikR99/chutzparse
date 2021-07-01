// +build wasm,electron

package model

import (
	"github.com/gontikr99/chutzparse/internal/eqlog"
	"github.com/gontikr99/chutzparse/internal/model/parsedefs"
	"github.com/gontikr99/chutzparse/internal/presenter"
)

func listenForHits() {
	eqlog.RegisterLogsListener(func(entries []*eqlog.LogEntry) {
		for _, entry := range entries {
			if dmgEntry, ok := entry.Meaning.(*eqlog.DamageLog); ok {
				if dmgEntry.Source == entry.Character && dmgEntry.Target == entry.Character {
					// self-inflicted damage
					presenter.BroadcastHitEvent(presenter.ChannelHitTop, &presenter.HitEvent{
						Text:  parsedefs.RenderAmount(float64(dmgEntry.Amount)),
						Color: "red",
						Big:   dmgEntry.Flag&eqlog.CriticalFlag != 0,
					})
				} else if dmgEntry.Source == entry.Character {
					// outgoing damage
					if dmgEntry.SpellName != "" {
						presenter.BroadcastHitEvent(presenter.ChannelHitTop, &presenter.HitEvent{
							Text:  parsedefs.RenderAmount(float64(dmgEntry.Amount)),
							Color: "yellow",
							Big:   dmgEntry.Flag&eqlog.CriticalFlag != 0,
						})
					} else if dmgEntry.Flag&eqlog.RiposteFlag != 0 {
						presenter.BroadcastHitEvent(presenter.ChannelHitTop, &presenter.HitEvent{
							Text:  parsedefs.RenderAmount(float64(dmgEntry.Amount)),
							Color: parsedefs.ColorPastelRed,
							Big:   dmgEntry.Flag&eqlog.CriticalFlag != 0,
						})
					} else if dmgEntry.Element == eqlog.PhysicalDamage {
						presenter.BroadcastHitEvent(presenter.ChannelHitTop, &presenter.HitEvent{
							Text:  parsedefs.RenderAmount(float64(dmgEntry.Amount)),
							Color: "white",
							Big:   dmgEntry.Flag&eqlog.CriticalFlag != 0,
						})
					} else {
						presenter.BroadcastHitEvent(presenter.ChannelHitTop, &presenter.HitEvent{
							Text:  parsedefs.RenderAmount(float64(dmgEntry.Amount)),
							Color: "gray",
							Big:   dmgEntry.Flag&eqlog.CriticalFlag != 0,
						})
					}
				} else if dmgEntry.Target == entry.Character {
					// incoming damage
					presenter.BroadcastHitEvent(presenter.ChannelHitBottom, &presenter.HitEvent{
						Text:  parsedefs.RenderAmount(float64(dmgEntry.Amount)),
						Color: "red",
						Big:   dmgEntry.Flag&eqlog.CriticalFlag != 0,
					})
				}
			} else if healEntry, ok := entry.Meaning.(*eqlog.HealLog); ok {
				if healEntry.Source == entry.Character && healEntry.Target == entry.Character {
					// self-healing
					presenter.BroadcastHitEvent(presenter.ChannelHitTop, &presenter.HitEvent{
						Text:  parsedefs.RenderAmount(float64(healEntry.Actual)),
						Color: parsedefs.ColorLimeGreen,
						Big:   healEntry.Flag&eqlog.CriticalFlag != 0,
					})
				} else if healEntry.Source == entry.Character {
					// outgoing healing
					presenter.BroadcastHitEvent(presenter.ChannelHitTop, &presenter.HitEvent{
						Text:  parsedefs.RenderAmount(float64(healEntry.Actual)),
						Color: "cyan",
						Big:   healEntry.Flag&eqlog.CriticalFlag != 0,
					})
				} else if healEntry.Target == entry.Character {
					// incoming healing
					presenter.BroadcastHitEvent(presenter.ChannelHitBottom, &presenter.HitEvent{
						Text:  parsedefs.RenderAmount(float64(healEntry.Actual)),
						Color: parsedefs.ColorLimeGreen,
						Big:   healEntry.Flag&eqlog.CriticalFlag != 0,
					})
				}
			}
		}
	})
}