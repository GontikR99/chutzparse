//go:build wasm && electron
// +build wasm,electron

package eqspec

import "github.com/gontikr99/chutzparse/pkg/multipattern"

type DeathLog struct {
	Source string
	Target string
}

func (d *DeathLog) Visit(handler ParsedLogHandler) interface{} {
	return handler.OnDeath(d)
}

func handleDeath(mp *multipattern.Multipattern) *multipattern.Multipattern {
	return commonSubpatterns(mp).
		On("(.*) died[.]", func(parts []string, _ interface{}) interface{} {
			return &DeathLog{
				Source: UnspecifiedName,
				Target: normalizeName(parts[1]),
			}
		}).
		On("(.*) has been slain by (.*)!", func(parts []string, _ interface{}) interface{} {
			return &DeathLog{
				Source: normalizeName(parts[2]),
				Target: normalizeName(parts[1]),
			}
		}).
		On("You have slain (.*)!", func(parts []string, _ interface{}) interface{} {
			return &DeathLog{
				Source: normalizeName("You"),
				Target: normalizeName(parts[1]),
			}
		}).
		On("You have been slain by (.*)!", func(parts []string, _ interface{}) interface{} {
			return &DeathLog{
				Source: normalizeName(parts[1]),
				Target: normalizeName("You"),
			}
		})
}
