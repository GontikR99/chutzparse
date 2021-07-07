// +build wasm,electron

package damage

import (
	"github.com/gontikr99/chutzparse/internal/eqlog"
	"github.com/gontikr99/chutzparse/internal/model/fight"
)

func (r *Report) Offer(entry *eqlog.LogEntry, epoch int) fight.FightReport {
	r.LastCharName = entry.Character
	if dmg, ok := entry.Meaning.(*eqlog.DamageLog); ok {
		if dmg.Target != r.Target && dmg.Target != r.Target+"`s pet" && dmg.Target != r.Target+"`s warder" {
			return r
		}
		contrib := r.ContributionOf(dmg.Source)
		contrib.TotalDamage += dmg.Amount
		cat := contrib.CategoryOf(dmg.DisplayCategory())
		cat.TotalDamage += dmg.Amount
		cat.Success++
	}
	return r
}
