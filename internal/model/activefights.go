// +build wasm,electron

package model

import (
	"bytes"
	"encoding/gob"
	"github.com/gontikr99/chutzparse/internal/eqspec"
	iff2 "github.com/gontikr99/chutzparse/internal/iff"
	"github.com/gontikr99/chutzparse/internal/model/fight"
	"github.com/gontikr99/chutzparse/internal/presenter"
	"github.com/gontikr99/chutzparse/internal/settings"
	"github.com/gontikr99/chutzparse/pkg/console"
	"github.com/gontikr99/chutzparse/pkg/electron/browserwindow"
	"sort"
	"strings"
	"time"
)

var activeFights = map[string]*fight.Fight{}
var activeUpdated bool

// nameReader collects the names of all characters involved in a log message
type nameReader struct {
	names map[string]struct{}
}

func (nr *nameReader) OnDamage(log *eqspec.DamageLog) interface{} {
	nr.storeName(log.Source)
	nr.storeName(log.Target)
	return nil
}
func (nr *nameReader) OnHeal(log *eqspec.HealLog) interface{} {
	nr.storeName(log.Source)
	nr.storeName(log.Target)
	return nil
}
func (nr *nameReader) OnDeath(log *eqspec.DeathLog) interface{} {
	nr.storeName(log.Source)
	nr.storeName(log.Target)
	return nil
}
func (nr *nameReader) OnChat(log *eqspec.ChatLog) interface{} { nr.storeName(log.Source); return nil }
func (nr *nameReader) OnZone(*eqspec.ZoneLog) interface{}     { return nil }
func (nr *nameReader) storeName(name string) {
	if name != eqspec.UnspecifiedName && name != "" {
		nr.names[name] = struct{}{}
	}
}

const inactivityTimeout = 12 * time.Second

func retireActiveFight(target string) {
	activeUpdated = true
	if fightData, present := activeFights[target]; present {
		fightData.Reports = fightData.Reports.Finalize(fightData)
		if fightData.Reports.Interesting() {
			buffer := &bytes.Buffer{}
			err := gob.NewEncoder(buffer).Encode(fightData)
			if err == nil {
				browserwindow.BroadcastChunked(fight.ChannelFinishedFights, buffer.Bytes())
			} else {
				console.Log(err)
			}
		}
	}
	delete(activeFights, target)
}

func listenForFights() {
	var fightIdGen int
	eqspec.RegisterLogsListener(func(entries []*eqspec.LogEntry) {
		// Any fight that hasn't been updated in a while gets retired
		now := time.Now()
		for _, fight := range activeFights {
			if now.Add(-inactivityTimeout).After(fight.StartTime) && now.Add(-inactivityTimeout).After(fight.LastActivity) {
				retireActiveFight(fight.Target)
			}
		}

		for _, entry := range entries {
			iff2.Update(entry)

			// If we just zoned, retire all active fights.
			if _, zoned := entry.Meaning.(*eqspec.ZoneLog); zoned {
				for target := range activeFights {
					retireActiveFight(target)
				}
				continue
			}

			// Create new fights as we encounter new NPCs
			if entry.Meaning != nil {
				nr := &nameReader{names: map[string]struct{}{}}
				entry.Meaning.Visit(nr)
				for name := range nr.names {
					if iff2.IsFoe(name) && iff2.GetOwner(name) == "" && !strings.HasSuffix(name, " pet") {
						if _, present := activeFights[name]; present {
							activeFights[name].LastActivity = time.Now()
						} else {
							activeFights[name] = &fight.Fight{
								Id:        fightIdGen,
								Target:    name,
								Reports:   fight.NewFightReports(name),
								StartTime: time.Now(),
							}
							fightIdGen++
							activeUpdated = true
						}
					}
				}
			}

			// Let every active fight hear about this log message
			for _, fight := range activeFights {
				for repName, report := range fight.Reports {
					fight.Reports[repName] = report.Offer(entry, fightIdGen)
					activeUpdated = true
				}
			}

			// Explicitly retire an active fight if its target dies
			if death, ok := entry.Meaning.(*eqspec.DeathLog); ok {
				if _, ok2 := activeFights[death.Target]; ok2 {
					retireActiveFight(death.Target)
				}
			}
		}
	})
}

type tsById []presenter.ThroughputState

func (t tsById) Len() int           { return len(t) }
func (t tsById) Less(i, j int) bool { return t[i].FightId < t[j].FightId }
func (t tsById) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }

func maintainThroughput() {
	go func() {
		for {
			<-time.After(333 * time.Millisecond) // damage meter update frequency
			if !activeUpdated {
				continue
			}
			activeUpdated = false
			var states []presenter.ThroughputState
			if value, present, err := settings.LookupSetting(settings.ShowMeters); err == nil && present && value == "true" {
				for _, fight := range activeFights {
					if !fight.Reports.Interesting() {
						continue
					}
					var top []presenter.ThroughputBar
					var bottom []presenter.ThroughputBar
					if dmgRep, present := fight.Reports["Damage"]; present {
						bottom = dmgRep.Throughput(fight)
					}
					if dmgRep, present := fight.Reports["Healing"]; present {
						top = dmgRep.Throughput(fight)
					}
					if top != nil || bottom != nil {
						states = append(states, presenter.ThroughputState{
							FightId:    fight.Id,
							TopBars:    top,
							BottomBars: bottom,
						})
					}
				}
				sort.Sort(tsById(states))
			}
			presenter.BroadcastThroughput(states)
		}
	}()
}

func init() {
	settings.DefaultSetting(settings.ShowMeters, "true")
}
