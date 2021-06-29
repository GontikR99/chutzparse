// +build wasm,web

package parsedefs

import (
	"github.com/vugu/vugu"
)

// FightReport records some specific aspect of a fight.
type FightReport interface {
	// Finalize the report, telling it that a fight underway has ended.
	Finalize() FightReport

	// Serialize this report for transmission
	Serialize() ([]byte, error)

	// Throughput generates a throughput chart as a summary from this report
	Throughput(fight *Fight) []ThroughputBar

	// Detail generates a detailed view of the information in this report
	Detail(fight *Fight) vugu.Builder
}
