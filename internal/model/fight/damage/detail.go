// +build wasm,web

package damage

import (
	"github.com/vugu/vugu"
)

type Detail struct {
	report *Report
}

func (r *Report) Detail() vugu.Builder {
	return &Detail{r}
}
