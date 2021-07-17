// +build wasm,electron

package boosts

import (
	"github.com/gontikr99/chutzparse/internal/eqspec"
	"github.com/gontikr99/chutzparse/internal/settings"
	"github.com/gontikr99/chutzparse/pkg/multipattern"
	"time"
)

// Update the set of boosts that are currently active by considering the given log message
func Update(entry *eqspec.LogEntry) {
	boostReader.Dispatch(entry.Message, entry)
	for _, cb := range activeBoosts {
		cb.maintain(entry.Timestamp)
	}
}

// Get the set of boosts that are currently active for the specified character
func Get(charName string) BoostSet {
	bs := BoostSet{}
	for bt, bm := range boostsForCharacter(charName) {
		for logId := range bm {
			bs.Add(bt, logId)
		}
	}
	return bs
}

var activeBoosts = map[string]characterBoosts{}

func boostsForCharacter(name string) characterBoosts {
	update, ok := activeBoosts[name]
	if !ok {
		update = characterBoosts{}
		activeBoosts[name] = update
	}
	return update
}

type characterBoosts map[BoostType]timedLogIds

func (cb characterBoosts) byType(bt BoostType) timedLogIds {
	update, ok := cb[bt]
	if !ok {
		update = timedLogIds{}
		cb[bt] = update
	}
	return update
}

// maintain removes expired boosts
func (cb characterBoosts) maintain(curTime time.Time) {
	for _, boostMap := range cb {
		for logId, expiration := range boostMap {
			if curTime.After(expiration) {
				delete(boostMap, logId)
			}
		}
	}
}

type timedLogIds map[int]time.Time

func addBoost(charName string, boostType BoostType, logId int, expiration time.Time) {
	boostsForCharacter(charName).byType(boostType)[logId] = expiration
}

var boostReader = multipattern.New().
	// Bard epics:
	// Prismatic Dragon Blade (Spirit of the Kin)
	On("You are filled with the spirit of the kin[.]", func(parts []string, context interface{}) interface{} {
		nbeText, _, _ := settings.LookupSetting(settings.NoteBardEpic)
		if nbeText=="true" {
			logEntry := context.(*eqspec.LogEntry)
			addBoost(logEntry.Character, BardEpic1_5, logEntry.Id, logEntry.Timestamp.Add(1*time.Minute))
		}
		return nil
	}).
	On("(.+) is filled with the spirit of the kin[.]", func(parts []string, context interface{}) interface{} {
		nbeText, _, _ := settings.LookupSetting(settings.NoteBardEpic)
		if nbeText=="true" {
			logEntry := context.(*eqspec.LogEntry)
			addBoost(logEntry.Character, BardEpic1_5, logEntry.Id, logEntry.Timestamp.Add(1*time.Minute))
		}
		return nil
	}).
	// Blade of Vesagran (Spirit of Vesagran)
	On("You are filled with the spirit of Vesagran[.]", func(parts []string, context interface{}) interface{} {
		nbeText, _, _ := settings.LookupSetting(settings.NoteBardEpic)
		if nbeText=="true" {
			logEntry := context.(*eqspec.LogEntry)
			addBoost(logEntry.Character, BardEpic2, logEntry.Id, logEntry.Timestamp.Add(1*time.Minute))
		}
		return nil
	}).
	On("(.+) is filled with the spirit of Vesagran[.]", func(parts []string, context interface{}) interface{} {
	nbeText, _, _ := settings.LookupSetting(settings.NoteBardEpic)
		if nbeText=="true" {
			logEntry := context.(*eqspec.LogEntry)
			addBoost(logEntry.Character, BardEpic2, logEntry.Id, logEntry.Timestamp.Add(1*time.Minute))
		}
		return nil
	}).

	// Shaman epics.  We assume shamans will have about 150% buff duration, so 1.5 minute uptimes.
	// Crafted Talisman of the Fates (Fateseer's Boon)
	On("You are blessed by the boon of the fateseer[.]", func(parts []string, context interface{}) interface{} {
		logEntry := context.(*eqspec.LogEntry)
		addBoost(logEntry.Character, ShamanEpic1_5, logEntry.Id, logEntry.Timestamp.Add(90*time.Second))
		return nil
	}).
	On("(.+) is blessed by the boon of the fateseer[.]", func(parts []string, context interface{}) interface{} {
		logEntry := context.(*eqspec.LogEntry)
		addBoost(parts[1], ShamanEpic1_5, logEntry.Id, logEntry.Timestamp.Add(90*time.Second))
		return nil
	}).
	// Blessed Spiritstaff of the Heyokah (Prophet's Gift of the Ruchu)
	On("You are blessed with the gift of the Ruchu[.]", func(parts []string, context interface{}) interface{} {
		logEntry := context.(*eqspec.LogEntry)
		addBoost(logEntry.Character, ShamanEpic2, logEntry.Id, logEntry.Timestamp.Add(90*time.Second))
		return nil
	}).
	On("(.+) is blessed with the gift of the Ruchu[.]", func(parts []string, context interface{}) interface{} {
		logEntry := context.(*eqspec.LogEntry)
		addBoost(parts[1], ShamanEpic2, logEntry.Id, logEntry.Timestamp.Add(90*time.Second))
		return nil
	}).

	// Intensity of the Resolute
	On("Your mind sharpens and strength flows into your body[.]", func(parts []string, context interface{}) interface{} {
		logEntry := context.(*eqspec.LogEntry)
		addBoost(logEntry.Character, IntensityResolute, logEntry.Id, logEntry.Timestamp.Add(60*time.Second))
		return nil
	}).
	On("(.+)'s mind sharpens and strength flows into ([a-z]+) body[.]", func(parts []string, context interface{}) interface{} {
		logEntry := context.(*eqspec.LogEntry)
		addBoost(parts[1], IntensityResolute, logEntry.Id, logEntry.Timestamp.Add(60*time.Second))
		return nil
	}).

	// Glyph of Destruction
	On("You activate your Glyph of Destruction[.]", func(parts []string, context interface{}) interface{} {
		logEntry := context.(*eqspec.LogEntry)
		addBoost(logEntry.Character, GlyphDestruction, logEntry.Id, logEntry.Timestamp.Add(120*time.Second))
		return nil
	}).
	On("(.+) is infused for destruction.[.]", func(parts []string, context interface{}) interface{} {
		logEntry := context.(*eqspec.LogEntry)
		addBoost(parts[1], GlyphDestruction, logEntry.Id, logEntry.Timestamp.Add(120*time.Second))
		return nil
	})

func init() {
	settings.DefaultSetting(settings.NoteBardEpic, "false")
}