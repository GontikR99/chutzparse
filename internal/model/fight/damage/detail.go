// +build wasm,web

package damage

import (
	"github.com/gontikr99/chutzparse/pkg/console"
	"github.com/vugu/vugu"
)

type Detail struct {
	report *Report
}

func (r *Report) Detail() vugu.Builder {
	console.Log("Generating detail")
	return &Detail{r}
}
