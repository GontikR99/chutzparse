package damage

import (
	"fmt"
	"github.com/gontikr99/chutzparse/internal/parse_model/parsedefs"
	"sort"
	"time"
)

type contribByDamageRev []*Contribution
func (c contribByDamageRev) Len() int {return len(c)}
func (c contribByDamageRev) Less(i, j int) bool {return c[i].DamageTotal() > c[j].DamageTotal()}
func (c contribByDamageRev) Swap(i, j int) {c[i], c[j] = c[j], c[i]}

func (r *Report) Throughput(fight *parsedefs.Fight) []parsedefs.ThroughputBar {
	duration := fight.LastActivity.Sub(fight.StartTime)
	if duration<0 {duration = time.Duration(1)}

	var contribs []*Contribution
	var totalDamage int64
	for _, contrib := range r.Contributions {
		contribs = append(contribs, contrib)
		totalDamage += contrib.DamageTotal()
	}

	durationSec := float64(duration)/float64(time.Second)
	dps := float64(totalDamage)/durationSec

	var bars []parsedefs.ThroughputBar
	bars = append (bars, parsedefs.ThroughputBar{
		CenterText: fmt.Sprintf("%s in %s", r.Target, parsedefs.RenderAmount(durationSec)),
		RightText:  fmt.Sprintf("%s = %s dps", parsedefs.RenderAmount(float64(totalDamage)), parsedefs.RenderAmount(dps)),
		Sectors:    []parsedefs.BarSector{{"green", 1.0}},
	})

	if len(contribs)==0 {
		return bars
	}
	sort.Sort(contribByDamageRev(contribs))
	maxDmg := contribs[0].DamageTotal()
	for i:=0;i<throughputBarCount && i<len(contribs);i++ {
		barColor := "blue"
		if contribs[i].Source == r.LastCharName {
			barColor = "red"
		}
		dps = float64(contribs[i].DamageTotal())/durationSec
		bars = append (bars, parsedefs.ThroughputBar{
			LeftText:   fmt.Sprintf("%d.", 1+i),
			CenterText: contribs[i].Source,
			RightText:  fmt.Sprintf("%s [%.3g%%] = %s dps",
				parsedefs.RenderAmount(float64(contribs[i].DamageTotal())),
				float64(100*contribs[i].DamageTotal())/float64(totalDamage),
				parsedefs.RenderAmount(dps)),
			Sectors:    []parsedefs.BarSector{{barColor, float64(contribs[i].DamageTotal())/float64(maxDmg)}},
		})
	}
	for i:=throughputBarCount; i<len(contribs); i++ {
		if contribs[i].Source != r.LastCharName {continue}
		dps = float64(contribs[i].DamageTotal())/durationSec
		bars = append (bars, parsedefs.ThroughputBar{
			LeftText:   fmt.Sprintf("%d.", 1+i),
			CenterText: contribs[i].Source,
			RightText:  fmt.Sprintf("%s [%.3g%%] = %s dps",
				parsedefs.RenderAmount(float64(contribs[i].DamageTotal())),
				float64(100*contribs[i].DamageTotal())/float64(totalDamage),
				parsedefs.RenderAmount(dps)),
			Sectors:    []parsedefs.BarSector{{"red", float64(contribs[i].DamageTotal())/float64(maxDmg)}},
		})
	}
	return bars
}

