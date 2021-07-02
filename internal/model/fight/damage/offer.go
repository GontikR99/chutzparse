// +build wasm,electron

package damage

import (
	"github.com/gontikr99/chutzparse/internal/eqlog"
	"github.com/gontikr99/chutzparse/internal/model/fight"
)

func (r *Report) Offer(entry *eqlog.LogEntry, epoch int) fight.FightReport {
	r.LastCharName = entry.Character
	dmg, ok := entry.Meaning.(*eqlog.DamageLog)
	if !ok {return r}
	if dmg.Target!=r.Target && dmg.Target!=r.Target+"`s pet" && dmg.Target!=r.Target+"`s warder" {
		return r
	}
	if _, ok := r.Contributions[dmg.Source]; !ok {
		r.Contributions[dmg.Source]=&Contribution{Source: dmg.Source}
	}
	r.Contributions[dmg.Source].TotalDamage+=dmg.Amount
	return r
}
