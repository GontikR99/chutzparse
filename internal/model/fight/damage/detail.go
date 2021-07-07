// +build wasm,web

package damage

import (
	"fmt"
	"github.com/gontikr99/chutzparse/internal/model/parsedefs"
	"github.com/vugu/vugu"
	"sort"
	"time"
)

type Detail struct {
	report *Report
	totalDamage int64
}

func (r *Report) Detail() vugu.Builder {
	result := &Detail{report: r}
	for _, contrib := range r.Contributions {
		result.totalDamage += contrib.DamageTotal()
	}
	return result
}

type detailRow struct {
	Source string
	Percent string
	Amount string
	DPS string
	damageDone int64
}

type drByDamageRev []*detailRow

func (d drByDamageRev) Len() int {return len(d)}
func (d drByDamageRev) Less(i, j int) bool {return d[i].damageDone > d[j].damageDone}
func (d drByDamageRev) Swap(i, j int) {d[i],d[j]=d[j],d[i]}

func (c *Detail) rows() []*detailRow {
	var result []*detailRow
	duration := c.report.EndTime.Sub(c.report.StartTime)
	for _, contrib := range c.report.Contributions {
		result = append(result, &detailRow{
			Source:     contrib.Source,
			Percent:    fmt.Sprintf("%.3g%%", float64(100*contrib.DamageTotal())/float64(c.totalDamage)),
			Amount:     fmt.Sprintf("%d", contrib.DamageTotal()),
			DPS:        parsedefs.RenderAmount(float64(contrib.DamageTotal())/float64(duration)*float64(time.Second)),
			damageDone: contrib.DamageTotal(),
		})
	}
	sort.Sort(drByDamageRev(result))
	return result
}