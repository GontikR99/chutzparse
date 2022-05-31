//go:build wasm && electron
// +build wasm,electron

package eqspec

import "github.com/gontikr99/chutzparse/pkg/multipattern"

type OutputfileComplete struct {
	Filename string
}

func (oc *OutputfileComplete) Visit(handler ParsedLogHandler) interface{} { return oc }

func handleDump(mp *multipattern.Multipattern) *multipattern.Multipattern {
	return commonSubpatterns(mp).
		On("Outputfile Complete: (.*[.]txt)", func(parts []string, _ interface{}) interface{} {
			return &OutputfileComplete{Filename: parts[1]}
		})
}
