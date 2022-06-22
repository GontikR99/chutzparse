//go:build wasm && electron
// +build wasm,electron

package eqspec

import (
	"github.com/gontikr99/chutzparse/pkg/multipattern"
	"strconv"
)

type RandomLog struct {
	Source string
	Lbound int32
	Ubound int32
	Value  int32
}

func (rl *RandomLog) Visit(handler ParsedLogHandler) interface{} { return rl }

func handleRandoms(mp *multipattern.Multipattern) *multipattern.Multipattern {
	return commonSubpatterns(mp).
		On("\\*\\*A Magic Die is rolled by (.+)[.] It could have been any number from (@num@) to (@num@),"+
			" but this time it turned up a (@num@)[.]", func(parts []string, context interface{}) interface{} {
			lb, _ := strconv.Atoi(parts[2])
			ub, _ := strconv.Atoi(parts[3])
			roll, _ := strconv.Atoi(parts[4])
			return &RandomLog{
				Source: normalizeName(parts[1]),
				Lbound: int32(lb),
				Ubound: int32(ub),
				Value:  int32(roll),
			}
		})
}
