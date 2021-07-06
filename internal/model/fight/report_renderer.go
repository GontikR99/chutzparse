// +build wasm,web

package fight

import (
	"github.com/vugu/vugu"
)

// FightReport records some specific aspect of a fight.
type FightReport interface {
	// Detail generates a detailed view of the information in this fight
	Detail() vugu.Builder
}
