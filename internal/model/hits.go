// +build wasm,electron

package model

import (
	"github.com/gontikr99/chutzparse/internal/eqspec"
	iff2 "github.com/gontikr99/chutzparse/internal/iff"
	"github.com/gontikr99/chutzparse/internal/model/parsedefs"
	"github.com/gontikr99/chutzparse/internal/presenter"
	"github.com/gontikr99/chutzparse/internal/settings"
)

func init() {
	settings.DefaultSetting(settings.ShowFlyingHits, "true")
}

func listenForHits() {
	eqspec.RegisterLogsListener(func(entries []*eqspec.LogEntry) {
		st, present, err := settings.LookupSetting(settings.ShowFlyingHits)
		if err==nil && present && st=="true" {
			for _, entry := range entries {
				if dmgEntry, ok := entry.Meaning.(*eqspec.DamageLog); ok {
					if dmgEntry.Source == entry.Character && dmgEntry.Target == entry.Character {
						// self-inflicted damage
						presenter.BroadcastHitEvent(presenter.ChannelHitTop, &presenter.HitEvent{
							Text:  parsedefs.FormatAmount(float64(dmgEntry.Amount)),
							Color: "red",
							Big:   dmgEntry.Flag&eqspec.CriticalFlag != 0,
						})
					} else if dmgEntry.Source == entry.Character {
						// outgoing damage
						if dmgEntry.SpellName != "" {
							presenter.BroadcastHitEvent(presenter.ChannelHitTop, &presenter.HitEvent{
								Text:  parsedefs.FormatAmount(float64(dmgEntry.Amount)),
								Color: "yellow",
								Big:   dmgEntry.Flag&eqspec.CriticalFlag != 0,
							})
						} else if dmgEntry.Flag&eqspec.RiposteFlag != 0 {
							presenter.BroadcastHitEvent(presenter.ChannelHitTop, &presenter.HitEvent{
								Text:  parsedefs.FormatAmount(float64(dmgEntry.Amount)),
								Color: parsedefs.ColorPastelRed,
								Big:   dmgEntry.Flag&eqspec.CriticalFlag != 0,
							})
						} else if dmgEntry.Element == eqspec.PhysicalDamage {
							presenter.BroadcastHitEvent(presenter.ChannelHitTop, &presenter.HitEvent{
								Text:  parsedefs.FormatAmount(float64(dmgEntry.Amount)),
								Color: "white",
								Big:   dmgEntry.Flag&eqspec.CriticalFlag != 0,
							})
						} else {
							presenter.BroadcastHitEvent(presenter.ChannelHitTop, &presenter.HitEvent{
								Text:  parsedefs.FormatAmount(float64(dmgEntry.Amount)),
								Color: "gray",
								Big:   dmgEntry.Flag&eqspec.CriticalFlag != 0,
							})
						}
					} else if dmgEntry.Target == entry.Character {
						// incoming damage
						presenter.BroadcastHitEvent(presenter.ChannelHitBottom, &presenter.HitEvent{
							Text:  parsedefs.FormatAmount(float64(dmgEntry.Amount)),
							Color: "red",
							Big:   dmgEntry.Flag&eqspec.CriticalFlag != 0,
						})
					} else if iff2.GetOwner(dmgEntry.Source) == entry.Character {
						// pet doing damage
						if dmgEntry.SpellName != "" {
							presenter.BroadcastHitEvent(presenter.ChannelHitTop, &presenter.HitEvent{
								Text:  parsedefs.FormatAmount(float64(dmgEntry.Amount)),
								Color: "khaki",
								Big:   dmgEntry.Flag&eqspec.CriticalFlag != 0,
							})
						} else if dmgEntry.Element == eqspec.PhysicalDamage {
							presenter.BroadcastHitEvent(presenter.ChannelHitTop, &presenter.HitEvent{
								Text:  parsedefs.FormatAmount(float64(dmgEntry.Amount)),
								Color: "bisque",
								Big:   dmgEntry.Flag&eqspec.CriticalFlag != 0,
							})
						} else {
							presenter.BroadcastHitEvent(presenter.ChannelHitTop, &presenter.HitEvent{
								Text:  parsedefs.FormatAmount(float64(dmgEntry.Amount)),
								Color: "gray",
								Big:   dmgEntry.Flag&eqspec.CriticalFlag != 0,
							})
						}
					} else if iff2.GetOwner(dmgEntry.Target) == entry.Character {
						// incoming damage
						presenter.BroadcastHitEvent(presenter.ChannelHitBottom, &presenter.HitEvent{
							Text:  parsedefs.FormatAmount(float64(dmgEntry.Amount)),
							Color: "indianred",
							Big:   dmgEntry.Flag&eqspec.CriticalFlag != 0,
						})
					}
				} else if healEntry, ok := entry.Meaning.(*eqspec.HealLog); ok {
					if healEntry.Source == entry.Character && healEntry.Target == entry.Character {
						// self-healing
						presenter.BroadcastHitEvent(presenter.ChannelHitTop, &presenter.HitEvent{
							Text:  parsedefs.FormatAmount(float64(healEntry.Actual)),
							Color: parsedefs.ColorLimeGreen,
							Big:   healEntry.Flag&eqspec.CriticalFlag != 0,
						})
					} else if healEntry.Source == entry.Character {
						// outgoing healing
						presenter.BroadcastHitEvent(presenter.ChannelHitTop, &presenter.HitEvent{
							Text:  parsedefs.FormatAmount(float64(healEntry.Actual)),
							Color: "cyan",
							Big:   healEntry.Flag&eqspec.CriticalFlag != 0,
						})
					} else if healEntry.Target == entry.Character {
						// incoming healing
						presenter.BroadcastHitEvent(presenter.ChannelHitBottom, &presenter.HitEvent{
							Text:  parsedefs.FormatAmount(float64(healEntry.Actual)),
							Color: parsedefs.ColorLimeGreen,
							Big:   healEntry.Flag&eqspec.CriticalFlag != 0,
						})
					}
				}
			}
		}
	})
}
