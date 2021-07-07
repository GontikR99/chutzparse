// +build wasm,web

package heal

import (
	"github.com/vugu/vugu"
)

func (r *Report) Detail() vugu.Builder {
	// FIXME: implement
	return &Detail{}
}

type Detail struct {
}