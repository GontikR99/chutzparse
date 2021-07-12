// +build wasm,electron

package heal

import (
	"github.com/gontikr99/chutzparse/internal/eqlog"
	iff2 "github.com/gontikr99/chutzparse/internal/iff"
	"github.com/gontikr99/chutzparse/internal/model/fight"
)

func (r *Report) Offer(entry *eqlog.LogEntry, epoch int) fight.FightReport {
	if healEntry, ok := entry.Meaning.(*eqlog.HealLog); ok {
		if !iff2.IsFoe(healEntry.Target) {
			contrib := r.ContributionOf(healEntry.Source)
			cat := contrib.CategoryOf(healEntry.DisplayCategory())
			epochStat := cat.EpochOf(epoch)

			contrib.TotalHealed += healEntry.Actual
			epochStat.TotalHealed += healEntry.Actual
			epochStat.Count++
		}
	}
	r.LastCharName = entry.Character
	return r
}
