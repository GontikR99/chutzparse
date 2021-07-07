// +build wasm,electron

package fight

import (
	"github.com/gontikr99/chutzparse/internal/eqlog"
	"github.com/gontikr99/chutzparse/internal/presenter"
)

// FightReport records some specific aspect of a fight.
// In the main process, we don't require that reports provide a detail view, but instead they must
// provide an Offer method.
type FightReport interface {
	// Offer this fight a new piece of information, generating an updated fight
	// here, `entry` represents the new information, while `epoch` is a value that
	// changes each time a new active fight is recorded
	Offer(entry *eqlog.LogEntry, epoch int) FightReport

	// Finalize the fight, telling it that a fight underway has ended.
	Finalize(fight *Fight) FightReport

	// Throughput generates a throughput chart as a summary from this fight
	Throughput(fight *Fight) []presenter.ThroughputBar

	// Interesting tells us whether this report has enough information to bother reporting the fight
	Interesting() bool
}

func (s FightReportSet) Interesting() bool {
	for _, report := range s {
		if report.Interesting() {
			return true
		}
	}
	return false
}

// Finalize an entire report set.
func (s FightReportSet) Finalize(fight *Fight) FightReportSet {
	for k, v := range s {
		s[k] = v.Finalize(fight)
	}
	return s
}
