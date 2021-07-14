// +build wasm,web

package tanking

import (
	"fmt"
	"github.com/gontikr99/chutzparse/internal/model/parsedefs"
	"github.com/gontikr99/chutzparse/internal/presenter"
	"github.com/vugu/vugu"
	"strconv"
	"strings"
	"time"
)

func (r *Report) Summarize() string {
	duration := r.ActivitySet.TotalDuration()
	if duration < time.Second {
		duration = time.Second
	}
	durationSec := float64(duration) / float64(time.Second)

	sb := &strings.Builder{}

	sb.WriteString(fmt.Sprintf("ChutzTank: %s [%s dealt] in %ss: ", r.Source,
		parsedefs.FormatAmount(float64(r.TotalDamage())),
		parsedefs.FormatAmount(durationSec)))

	contribs := r.SortedContributors()
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
		sb.WriteString(fmt.Sprintf("%d. %s [%s=%s]", 1+i,
			contribs[i].Target,
			parsedefs.FormatAmount(float64(contribs[i].TotalDamage)),
			parsedefs.FormatPercent(float64(contribs[i].TotalDamage)/float64(r.TotalDamage())),
		))
	}
	return sb.String()

}

func (r *Report) Detail() vugu.Builder {
	return &Detail{report: r}
}

type Detail struct {
	report *Report
}

type displayRow struct {
	Rank string
	Target string
	Percent string
	Total string
}

func (c *Detail) rows() []*displayRow {
	totalDamage := c.report.TotalDamage()
	var result []*displayRow
	for idx, contrib :=range c.report.SortedContributors() {
		if contrib.TotalDamage!=0 {
			result = append(result, &displayRow{
				Rank:    strconv.FormatInt(1+int64(idx), 10),
				Target:  contrib.Target,
				Percent: parsedefs.FormatPercent(float64(contrib.TotalDamage)/float64(totalDamage)),
				Total:   parsedefs.FormatAmount(float64(contrib.TotalDamage)),
			})
		}
	}
	return result
}