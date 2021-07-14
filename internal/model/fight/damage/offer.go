// +build wasm,electron

package damage

import (
	"github.com/gontikr99/chutzparse/internal/eqspec"
	"github.com/gontikr99/chutzparse/internal/iff"
	"github.com/gontikr99/chutzparse/internal/model/boosts"
	"github.com/gontikr99/chutzparse/internal/model/fight"
)

func (r *Report) Offer(entry *eqspec.LogEntry, epoch int) fight.FightReport {
	r.LastCharName = entry.Character
	if dmg, ok := entry.Meaning.(*eqspec.DamageLog); ok {
		if dmg.Target != r.Target && iff.GetOwner(dmg.Target) != r.Target {
			return r
		}
		contrib := r.ContributionOf(dmg.Source)
		contrib.Boosts.AddAll(boosts.Get(dmg.Source))
		contrib.TotalDamage += dmg.Amount
		cat := contrib.CategoryOf(dmg.DisplayCategory())
		cat.TotalDamage += dmg.Amount
		cat.Success++
	}
	if death, ok := entry.Meaning.(*eqspec.DeathLog); ok {
		contrib := r.ContributionOf(death.Target)
		contrib.Deaths[entry.Id] = struct{}{}
	}
	return r
}
