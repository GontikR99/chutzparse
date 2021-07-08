// +build wasm,electron

package eqlog

import (
	"fmt"
	"github.com/gontikr99/chutzparse/pkg/multipattern"
	"strings"
)

type HealLog struct {
	Source string
	Target string
	Total  int64
	Actual int64

	Flag      HitFlag
	SpellName string
}

func (h *HealLog) Visit(handler ParsedLogHandler) interface{} {
	return handler.OnHeal(h)
}

func (h *HealLog) String() string {
	return fmt.Sprintf("Heal: %s -> %s Total: %d Actual: %d [%s] %s",
		h.Source, h.Target, h.Total, h.Actual, h.Flag, h.SpellName,
	)
}

func (h *HealLog) DisplayCategory() string {
	if h.SpellName != "" {
		return h.SpellName
	} else {
		return "unspecified"
	}
}

func handleHeal(mp *multipattern.Multipattern) *multipattern.Multipattern {
	return commonSubpatterns(mp).
		On("(.+) has been healed over time for (@num@) (?:\\((@num@)\\) )?hit points(?: by (.+))?[.](@hflag@)?", func(parts []string) interface{} {
			var actual int64
			var total int64
			actual = amount(parts[2])
			if parts[3] == "" {
				total = actual
			} else {
				total = amount(parts[3])
			}
			res := &HealLog{
				Source:    UnspecifiedName,
				Target:    normalizeName(parts[1]),
				Total:     total,
				Actual:    actual,
				Flag:      hitFlags(parts[5]),
				SpellName: parts[4],
			}

			if res.Target == "Yourself" || res.Target == "Himself" || res.Target == "Herself" || res.Target == "Itself" {
				res.Target = res.Source
			}
			return res
		}).
		On("(.*) healed (.*) for (@num@) (?:\\((@num@)\\) )?hit points(?: by (.*))?[.](@hflag@)?", func(parts []string) interface{} {
			if strings.HasSuffix(parts[2], " over time") {
				parts[2] = parts[2][:len(parts[2])-10]
			}
			var actual int64
			var total int64
			actual = amount(parts[3])
			if parts[4] == "" {
				total = actual
			} else {
				total = amount(parts[4])
			}
			res := &HealLog{
				Source:    normalizeName(parts[1]),
				Target:    normalizeName(parts[2]),
				Total:     total,
				Actual:    actual,
				Flag:      hitFlags(parts[6]),
				SpellName: parts[5],
			}

			if res.Target == "Yourself" || res.Target == "Himself" || res.Target == "Herself" || res.Target == "Itself" {
				res.Target = res.Source
			}
			return res
		}).
		On("The gods have healed (.*) for (@num@) points of damage[.]", func(parts []string) interface{} {
			return &HealLog{
				Source:    UnspecifiedName,
				Target:    parts[1],
				Total:     amount(parts[2]),
				Actual:    amount(parts[2]),
				Flag:      0,
				SpellName: "Divine Intervention",
			}
		})
}
