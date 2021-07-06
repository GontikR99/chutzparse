// +build wasm,web

package fight

import (
	"github.com/gontikr99/chutzparse/internal/model/parsedefs"
	"github.com/gontikr99/chutzparse/internal/model/presenter"
	"github.com/vugu/vugu"
)

// FightReport records some specific aspect of a fight.
type FightReport interface {
	// Finalize the fight, telling it that a fight underway has ended.
	Finalize() FightReport

	// Throughput generates a throughput chart as a summary from this fight
	Throughput(fight *parsedefs.Fight) []presenter.ThroughputBar

	// Detail generates a detailed view of the information in this fight
	Detail(fight *parsedefs.Fight) vugu.Builder
}
