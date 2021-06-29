// +build wasm,electron

package parse_model

import (
	"github.com/gontikr99/chutzparse/internal/eqlog"
	"github.com/gontikr99/chutzparse/internal/parse_model/parsecomms"
	"github.com/gontikr99/chutzparse/internal/parse_model/parsedefs"
	"github.com/gontikr99/chutzparse/pkg/console"
	"strings"
	"time"
	"unicode"
)

var activeFights=map[string]*parsedefs.Fight{}

// npcReader determines if an event involves an NPC, and stores the last NPC involved
type npcReader struct {
	npcs map[string]struct{}
}

func (nr *npcReader) OnDamage(log *eqlog.DamageLog) interface{} {nr.storeNPC(log.Source); nr.storeNPC(log.Target); return nil}
func (nr *npcReader) OnHeal(log *eqlog.HealLog) interface{} {nr.storeNPC(log.Source); nr.storeNPC(log.Target); return nil}
func (nr *npcReader) OnDeath(log *eqlog.DeathLog) interface{} {nr.storeNPC(log.Source); nr.storeNPC(log.Target); return nil}
func (nr *npcReader) OnZone(log *eqlog.ZoneLog) interface{} {return nil}

func (nr *npcReader) storeNPC(name string) {
	// Damage done by "Pain and Suffering" isn't all that useful
	if name==eqlog.UnspecifiedName {return}

	// Corpses don't count at all
	if strings.HasSuffix(name, "`s corpse") {
		return
	}

	// Warders and pets don't count, unless they're NPC warders/pets
	if strings.HasSuffix(name, "`s warder") {
		nr.storeNPC(name[:len(name)-9])
		return
	}
	if strings.HasSuffix(name, "`s pet") {
		nr.storeNPC(name[:len(name)-6])
		return
	}

	// all PC names start with a capital letter
	if name!="" && unicode.IsLower(rune(name[0])) {
		nr.npcs[name]=struct{}{}
		return
	}

	// Heuristic: all NPC names contain spaces
	// FIXME: maybe enumerate the few NPCs whose names don't contain spaces?
	if !strings.ContainsRune(name, ' ') {
		return
	}
	nr.npcs[name]=struct{}{}
}

const inactivityTimeout=12*time.Second

func retireActiveFight(target string) {
	console.Logf("Retiring %s", target)
	delete(activeFights, target)
}

func listenForFights() {
	var fightIdGen int
	eqlog.RegisterLogsListener(func(entries []*eqlog.LogEntry) {
		// Any fight that hasn't been updated in a while gets retired
		now := time.Now()
		for _, fight := range activeFights {
			if now.Add(-inactivityTimeout).After(fight.StartTime) && now.Add(-inactivityTimeout).After(fight.LastActivity) {
				retireActiveFight(fight.Target)
			}
		}

		for _, entry := range entries {
			// If we just zoned, retire all active fights.
			if _, zoned := entry.Meaning.(*eqlog.ZoneLog); zoned {
				for target, _ := range activeFights {
					retireActiveFight(target)
				}
				continue
			}

			// Create new fights as we encounter new NPCs
			if entry.Meaning!=nil {
				nr := &npcReader{npcs: map[string]struct{}{}}
				entry.Meaning.Visit(nr)
				for npc, _ := range nr.npcs {
					if _, present := activeFights[npc]; present {
						activeFights[npc].LastActivity = time.Now()
					} else {
						activeFights[npc] = &parsedefs.Fight{
							Id:        fightIdGen,
							Target:    npc,
							Reports:   parsedefs.NewFightReports(npc),
							StartTime: time.Now(),
						}
						fightIdGen++
					}
				}
			}

			// Let every active fight hear about this log message
			for _, fight := range activeFights {
				for repName, report := range fight.Reports {
					fight.Reports[repName]=report.Offer(entry, fightIdGen)
				}
			}

			// Explicitly retire an active fight if its target dies
			if death, ok := entry.Meaning.(*eqlog.DeathLog); ok {
				if _, ok2 := activeFights[death.Target]; ok2 {
					retireActiveFight(death.Target)
				}
			}
		}
	})
}

type tsById []parsedefs.ThroughputState

func (t tsById) Len() int {return len(t)}
func (t tsById) Less(i, j int) bool {return t[i].FightId < t[j].FightId}
func (t tsById) Swap(i, j int) {t[i], t[j] = t[j], t[i]}

func maintainThroughput() {
	go func() {
		for {
			<-time.After(1*time.Second)
			var states []parsedefs.ThroughputState
			for _, fight := range activeFights {
				if dmgRep, present := fight.Reports["Damage"]; present {
					states = append(states, parsedefs.ThroughputState{
						FightId:    fight.Id,
						TopBars:    nil,
						BottomBars: dmgRep.Throughput(fight),
					})
				}
			}
			parsecomms.BroadcastThroughput(states)
		}
	}()
}