// +build wasm,electron

package parsedefs// FightReport records some specific aspect of a fight.
import (
	"github.com/gontikr99/chutzparse/internal/eqlog"
)

type FightReport interface {
	// Offer this report a new piece of information, generating an updated report
	// here, `entry` represents the new information, while `epoch` is a value that
	// changes each time a new active fight is recorded
	Offer(entry *eqlog.LogEntry, epoch int) FightReport

	// Finalize the report, telling it that a fight underway has ended.
	Finalize() FightReport

	// Serialize this report for transmission
	Serialize() ([]byte, error)

	// Throughput generates a throughput chart as a summary from this report
	Throughput(fight *Fight) []ThroughputBar
}
