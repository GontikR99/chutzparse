//go:build wasm && web
// +build wasm,web

package damage

import (
	"fmt"
	"github.com/gontikr99/chutzparse/internal/model/parsedefs"
	"github.com/gontikr99/chutzparse/internal/presenter"
	"strings"
	"time"
)

func (r *Report) Summarize() string {
	agRep := r.Aggregate()
	duration := r.ActivitySet.TotalDuration()
	if duration < time.Second {
		duration = time.Second
	}
	durationSec := float64(duration) / float64(time.Second)

	sb := &strings.Builder{}

	sb.WriteString(fmt.Sprintf("ChutzDmg: %s [%s HP] in %ss: ", r.Target,
		parsedefs.FormatAmount(float64(agRep.TotalDamage)),
		parsedefs.FormatAmount(durationSec)))

	contribs := agRep.SortedContributors()
	if len(contribs) == 0 {
		return ""
	}
	needSep := false
	for i := 0; i < len(contribs) && i < presenter.ThroughputBarCount; i++ {
		if needSep {
			sb.WriteString("; ")
		} else {
			needSep = true
		}
		flags := contribs[i].Flags
		if flags != "" {
			flags = " {" + flags + "}"
		}
		sb.WriteString(fmt.Sprintf("%d. %s%s %s DPS [%s]", 1+i,
			contribs[i].AttributedSource,
			flags,
			parsedefs.FormatAmount(float64(contribs[i].TotalDamage)/float64(durationSec)),
			parsedefs.FormatPercent(float64(contribs[i].TotalDamage)/float64(agRep.TotalDamage)),
		))
	}
	return sb.String()
}
