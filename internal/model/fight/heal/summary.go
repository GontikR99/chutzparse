// +build wasm,web

package heal

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
		duration=time.Second
	}
	durationSec := float64(duration)/float64(time.Second)

	sb := &strings.Builder{}

	sb.WriteString(fmt.Sprintf("ChutzHeals: %s in %ss: ", r.Belligerent, parsedefs.FormatAmount(durationSec)))

	contribs := agRep.SortedContributors()
	if len(contribs)==0 {
		return ""
	}
	needSep := false
	for i:=0; i<len(contribs) && i<presenter.ThroughputBarCount ; i++ {
		if needSep {
			sb.WriteString("; ")
		} else {
			needSep=true
		}
		sb.WriteString(fmt.Sprintf("%d. %s %s HPS [%s]", 1+i,
			contribs[i].AttributedSource,
			parsedefs.FormatAmount(float64(contribs[i].TotalHealed)/float64(durationSec)),
			parsedefs.FormatPercent(float64(contribs[i].TotalHealed)/float64(agRep.TotalHealed)),
		))
	}
	return sb.String()
}