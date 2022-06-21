//go:build wasm && electron
// +build wasm,electron

package model

import (
	"github.com/gontikr99/chutzparse/internal/eqspec"
	"github.com/gontikr99/chutzparse/internal/iff"
	"github.com/gontikr99/chutzparse/internal/model/parsedefs"
	"github.com/gontikr99/chutzparse/internal/presenter"
	"github.com/gontikr99/chutzparse/internal/settings"
)

var defaultColors = map[string]string{
	settings.HitColorSpellDamage:   "#ffff00",
	settings.HitColorRiposteDamage: "#faa0a0",
	settings.HitColorMeleeDamage:   "#ffffff",
	settings.HitColorOtherDamage:   "#808080",
	settings.HitColorDamageTaken:   "#ff0000",

	settings.HitColorPetSpellDamage:   "#f0e68c",
	settings.HitColorPetRiposteDamage: "#ffe4c4",
	settings.HitColorPetMeleeDamage:   "#ffe4c4",
	settings.HitColorPetOtherDamage:   "#808080",
	settings.HitColorPetDamageTaken:   "#cd5c5c",

	settings.HitColorHealingReceived: "#32cd32",
	settings.HitColorHealingDone:     "#00ffff",

	settings.HitColorPetHealingReceived: "#32cd32",
	settings.HitColorPetHealingDone:     "#00ffff",
}

func init() {
	settings.DefaultSetting(settings.ShowFlyingHits, "true")

	for key, color := range defaultColors {
		settings.DefaultSetting(key, color)
	}
}

func colorOf(name string) string {
	if value, present, err := settings.LookupSetting(name); err == nil && present {
		return value
	} else {
		return "#663399"
	}
}

func listenForHits() {
	eqspec.RegisterLogsListener(func(entries []*eqspec.LogEntry) {
		st, present, err := settings.LookupSetting(settings.ShowFlyingHits)
		if err == nil && present && st == "true" {
			for _, entry := range entries {
				if dmgEntry, ok := entry.Meaning.(*eqspec.DamageLog); ok {
					if dmgEntry.Source == entry.Character && dmgEntry.Target == entry.Character {
						// self-inflicted damage
						presenter.BroadcastHitEvent(presenter.ChannelHitTop, &presenter.HitEvent{
							Text:  parsedefs.FormatAmount(float64(dmgEntry.Amount)),
							Color: colorOf(settings.HitColorDamageTaken),
							Big:   dmgEntry.Flag&eqspec.CriticalFlag != 0,
						})
					} else if dmgEntry.Source == entry.Character {
						// outgoing damage
						if dmgEntry.SpellName != "" {
							presenter.BroadcastHitEvent(presenter.ChannelHitTop, &presenter.HitEvent{
								Text:  parsedefs.FormatAmount(float64(dmgEntry.Amount)),
								Color: colorOf(settings.HitColorSpellDamage),
								Big:   dmgEntry.Flag&eqspec.CriticalFlag != 0,
							})
						} else if dmgEntry.Flag&eqspec.RiposteFlag != 0 {
							presenter.BroadcastHitEvent(presenter.ChannelHitTop, &presenter.HitEvent{
								Text:  parsedefs.FormatAmount(float64(dmgEntry.Amount)),
								Color: colorOf(settings.HitColorRiposteDamage),
								Big:   dmgEntry.Flag&eqspec.CriticalFlag != 0,
							})
						} else if dmgEntry.Element == eqspec.PhysicalDamage {
							presenter.BroadcastHitEvent(presenter.ChannelHitTop, &presenter.HitEvent{
								Text:  parsedefs.FormatAmount(float64(dmgEntry.Amount)),
								Color: colorOf(settings.HitColorMeleeDamage),
								Big:   dmgEntry.Flag&eqspec.CriticalFlag != 0,
							})
						} else {
							presenter.BroadcastHitEvent(presenter.ChannelHitTop, &presenter.HitEvent{
								Text:  parsedefs.FormatAmount(float64(dmgEntry.Amount)),
								Color: colorOf(settings.HitColorOtherDamage),
								Big:   dmgEntry.Flag&eqspec.CriticalFlag != 0,
							})
						}
					} else if dmgEntry.Target == entry.Character {
						// incoming damage
						presenter.BroadcastHitEvent(presenter.ChannelHitBottom, &presenter.HitEvent{
							Text:  parsedefs.FormatAmount(float64(dmgEntry.Amount)),
							Color: colorOf(settings.HitColorDamageTaken),
							Big:   dmgEntry.Flag&eqspec.CriticalFlag != 0,
						})
					} else if iff.GetOwner(dmgEntry.Source) == entry.Character {
						// pet doing damage
						if dmgEntry.SpellName != "" {
							presenter.BroadcastHitEvent(presenter.ChannelHitTop, &presenter.HitEvent{
								Text:  parsedefs.FormatAmount(float64(dmgEntry.Amount)),
								Color: colorOf(settings.HitColorPetSpellDamage),
								Big:   dmgEntry.Flag&eqspec.CriticalFlag != 0,
							})
						} else if dmgEntry.Flag&eqspec.RiposteFlag != 0 {
							presenter.BroadcastHitEvent(presenter.ChannelHitTop, &presenter.HitEvent{
								Text:  parsedefs.FormatAmount(float64(dmgEntry.Amount)),
								Color: colorOf(settings.HitColorPetRiposteDamage),
								Big:   dmgEntry.Flag&eqspec.CriticalFlag != 0,
							})
						} else if dmgEntry.Element == eqspec.PhysicalDamage {
							presenter.BroadcastHitEvent(presenter.ChannelHitTop, &presenter.HitEvent{
								Text:  parsedefs.FormatAmount(float64(dmgEntry.Amount)),
								Color: colorOf(settings.HitColorPetMeleeDamage),
								Big:   dmgEntry.Flag&eqspec.CriticalFlag != 0,
							})
						} else {
							presenter.BroadcastHitEvent(presenter.ChannelHitTop, &presenter.HitEvent{
								Text:  parsedefs.FormatAmount(float64(dmgEntry.Amount)),
								Color: colorOf(settings.HitColorPetOtherDamage),
								Big:   dmgEntry.Flag&eqspec.CriticalFlag != 0,
							})
						}
					} else if iff.GetOwner(dmgEntry.Target) == entry.Character {
						// incoming damage
						presenter.BroadcastHitEvent(presenter.ChannelHitBottom, &presenter.HitEvent{
							Text:  parsedefs.FormatAmount(float64(dmgEntry.Amount)),
							Color: colorOf(settings.HitColorPetDamageTaken),
							Big:   dmgEntry.Flag&eqspec.CriticalFlag != 0,
						})
					}
				} else if healEntry, ok := entry.Meaning.(*eqspec.HealLog); ok {
					if healEntry.Source == entry.Character && healEntry.Target == entry.Character {
						// self-healing
						presenter.BroadcastHitEvent(presenter.ChannelHitTop, &presenter.HitEvent{
							Text:  parsedefs.FormatAmount(float64(healEntry.Actual)),
							Color: colorOf(settings.HitColorHealingReceived),
							Big:   healEntry.Flag&eqspec.CriticalFlag != 0,
						})
					} else if healEntry.Source == entry.Character {
						// outgoing healing
						presenter.BroadcastHitEvent(presenter.ChannelHitTop, &presenter.HitEvent{
							Text:  parsedefs.FormatAmount(float64(healEntry.Actual)),
							Color: colorOf(settings.HitColorHealingDone),
							Big:   healEntry.Flag&eqspec.CriticalFlag != 0,
						})
					} else if healEntry.Target == entry.Character {
						// incoming healing
						presenter.BroadcastHitEvent(presenter.ChannelHitBottom, &presenter.HitEvent{
							Text:  parsedefs.FormatAmount(float64(healEntry.Actual)),
							Color: colorOf(settings.HitColorHealingReceived),
							Big:   healEntry.Flag&eqspec.CriticalFlag != 0,
						})
					} else if iff.GetOwner(healEntry.Source) == entry.Character {
						// pet doing healing
						if healEntry.Target == healEntry.Source {
							// pet healing itself
							presenter.BroadcastHitEvent(presenter.ChannelHitTop, &presenter.HitEvent{
								Text:  parsedefs.FormatAmount(float64(healEntry.Actual)),
								Color: colorOf(settings.HitColorPetHealingReceived),
								Big:   healEntry.Flag&eqspec.CriticalFlag != 0,
							})
						} else {
							// pet outgoing healing
							presenter.BroadcastHitEvent(presenter.ChannelHitTop, &presenter.HitEvent{
								Text:  parsedefs.FormatAmount(float64(healEntry.Actual)),
								Color: colorOf(settings.HitColorPetHealingDone),
								Big:   healEntry.Flag&eqspec.CriticalFlag != 0,
							})
						}
					} else if iff.GetOwner(healEntry.Target) == entry.Character {
						// pet being healed
						presenter.BroadcastHitEvent(presenter.ChannelHitBottom, &presenter.HitEvent{
							Text:  parsedefs.FormatAmount(float64(healEntry.Actual)),
							Color: colorOf(settings.HitColorPetHealingReceived),
							Big:   healEntry.Flag&eqspec.CriticalFlag != 0,
						})
					}
				}
			}
		}
	})
}
