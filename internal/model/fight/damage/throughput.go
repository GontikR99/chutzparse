// +build wasm,electron

package damage

import (
	"fmt"
	"github.com/gontikr99/chutzparse/internal/model/fight"
	"github.com/gontikr99/chutzparse/internal/model/parsedefs"
	"github.com/gontikr99/chutzparse/internal/presenter"
	"time"
)

func toThroughputBar(c *aggregateContributor, index int, totalDamage int64, maxDmg int64, durationSec float64, colors []string) presenter.ThroughputBar {
	dps := float64(c.TotalDamage) / durationSec
	result := presenter.ThroughputBar{
		LeftText:   fmt.Sprintf("%d. %s", index+1, c.DisplayName()),
		CenterText: "",
		RightText: fmt.Sprintf("%s [%.3g%%] = %s dps",
			parsedefs.RenderAmount(float64(c.TotalDamage)),
			float64(100*c.TotalDamage)/float64(totalDamage),
			parsedefs.RenderAmount(dps)),
		Sectors: nil,
	}
	for idx, ctb := range c.RawContributions {
		result.Sectors = append(result.Sectors, presenter.BarSector{
			Color:   colors[idx%len(colors)],
			Portion: float64(ctb.TotalDamage) / float64(maxDmg),
		})
	}
	return result
}

const throughputBarCount = 10
var otherColors = []string{"blue", "lightblue"}
var selfColors = []string{"red", "darksalmon"}

func (r *Report) Throughput(fight *fight.Fight) []presenter.ThroughputBar {
	duration := fight.LastActivity.Sub(fight.StartTime)
	if duration <= 1 {
		duration = 1 * time.Second
	}
	durationSec := float64(duration) / float64(time.Second)

	aggRep := r.Aggregate()
	dps := float64(aggRep.TotalDamage) / durationSec

	var bars []presenter.ThroughputBar
	bars = append(bars, presenter.ThroughputBar{
		LeftText:  fmt.Sprintf("[Damage] %s in %ss", r.Target, parsedefs.RenderAmount(durationSec)),
		RightText: fmt.Sprintf("%s = %s dps", parsedefs.RenderAmount(float64(aggRep.TotalDamage)), parsedefs.RenderAmount(dps)),
		Sectors:   []presenter.BarSector{{"dimgray", 1.0}},
	})

	if len(aggRep.Contributions) == 0 {
		return bars
	}
	agContribs := aggRep.SortedContributors()
	maxDmg := agContribs[0].TotalDamage

	for i:=0; i < throughputBarCount && i<len(agContribs); i++ {
		colors := otherColors
		if agContribs[i].AttributedSource == r.LastCharName {
			colors = selfColors
		}
		bars = append(bars, toThroughputBar(agContribs[i], i, aggRep.TotalDamage, maxDmg, durationSec, colors))
	}
	for i := throughputBarCount; i<len(agContribs); i++ {
		if agContribs[i].AttributedSource == r.LastCharName {
			bars = append(bars, toThroughputBar(agContribs[i], i, aggRep.TotalDamage, maxDmg, durationSec, selfColors))
		}
	}
	return bars
}
