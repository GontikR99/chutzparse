//go:build wasm && electron
// +build wasm,electron

package eqspec

import "github.com/gontikr99/chutzparse/pkg/multipattern"

type ZoneLog struct {
	ZoneFull string
}

func (z *ZoneLog) Visit(handler ParsedLogHandler) interface{} { return handler.OnZone(z) }

func handleZone(mp *multipattern.Multipattern) *multipattern.Multipattern {
	return commonSubpatterns(mp).
		On("You have entered (.+)[.]", func(parts []string, _ interface{}) interface{} {
			return &ZoneLog{
				ZoneFull: parts[1],
			}
		})
}
