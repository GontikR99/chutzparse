// +build wasm,electron

package heal

import (
	"github.com/gontikr99/chutzparse/internal/eqlog"
	"github.com/gontikr99/chutzparse/internal/model/fight"
)

func (r *Report) Offer(entry *eqlog.LogEntry, epoch int) fight.FightReport {
	if healEntry, ok := entry.Meaning.(*eqlog.HealLog); ok {
		var update *Contribution
		if update, ok = r.Contributions[healEntry.Source]; !ok {
			update = &Contribution{
				Source:      healEntry.Source,
				TotalHealed: 0,
				HealByEpoch: make(map[int]int64),
			}
			r.Contributions[healEntry.Source]=update
		}
		update.TotalHealed += healEntry.Actual
		if _, present := update.HealByEpoch[epoch]; !present {
			update.HealByEpoch[epoch]=0
		}
		update.HealByEpoch[epoch]+=healEntry.Actual
	}
	r.LastCharName = entry.Character
	return r
}
