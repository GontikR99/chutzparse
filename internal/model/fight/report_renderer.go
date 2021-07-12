// +build wasm,web

package fight

import (
	"github.com/vugu/vugu"
)

// FightReport records some specific aspect of a fight.
type FightReport interface {
	// Detail generates a detailed view of the information in this fight
	Detail() vugu.Builder

	// Create a string summary of this report, for pasting to a clipboard
	Summarize() string

	// Update a set with all possible player controlled participants
	Participants(p map[string]struct{})
}

func (f FightReportSet) Participants() map[string]struct{} {
	result := map[string]struct{}{}
	for _, report := range f {
		report.Participants(result)
	}
	return result
}
