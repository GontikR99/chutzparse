// +build wasm,electron

package heal

import (
	"github.com/gontikr99/chutzparse/internal/eqlog"
	"github.com/gontikr99/chutzparse/internal/model/fight"
	"github.com/gontikr99/chutzparse/internal/model/iff"
)

func (r *Report) Offer(entry *eqlog.LogEntry, epoch int) fight.FightReport {
	if healEntry, ok := entry.Meaning.(*eqlog.HealLog); ok {
		if !iff.IsFoe(healEntry.Target) {
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
