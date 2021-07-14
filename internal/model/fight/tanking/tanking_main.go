// +build wasm,electron

package tanking

import (
	"fmt"
	"github.com/gontikr99/chutzparse/internal/eqspec"
	"github.com/gontikr99/chutzparse/internal/model/fight"
	"github.com/gontikr99/chutzparse/internal/model/parsedefs"
	"github.com/gontikr99/chutzparse/internal/presenter"
	"time"
)

func (r *Report) Offer(entry *eqspec.LogEntry, epoch int) fight.FightReport {
	r.LastCharName = entry.Character
	if dmg, ok := entry.Meaning.(*eqspec.DamageLog); ok {
		if dmg.Source != r.Source {return r}
		contrib := r.ContributionOf(dmg.Target)
		contrib.TotalDamage += dmg.Amount
		contrib.Hits++
	}
	return r
}

var otherColor = "blue"
var selfColor = "darkviolet"

func (r *Report) Throughput(fight *fight.Fight) []presenter.ThroughputBar {
	duration := fight.LastActivity.Sub(fight.StartTime)
	if duration <= 1 {
		duration = 1 * time.Second
	}
	durationSec := float64(duration) / float64(time.Second)

	totalDamage := r.TotalDamage()
	dps := float64(totalDamage)/durationSec

	var bars []presenter.ThroughputBar
	bars = append(bars, presenter.ThroughputBar{
		LeftText:   fmt.Sprintf("[Tanking] %s in %ss", r.Source, parsedefs.FormatFixed(durationSec)),
		RightText:  fmt.Sprintf("%s = %s dps", parsedefs.FormatFixed(float64(totalDamage)), parsedefs.FormatFixed(dps)),
		RightStyle: presenter.MonoStyle,
		Sectors:    []presenter.BarSector{{"dimgray", 1.0}},
	})

	if len(r.Contributions)==0 {return bars}
	contribs := r.SortedContributors()
	for idx := 0; idx<presenter.ThroughputBarCount && idx<len(contribs); idx++ {
		color := otherColor
		contrib := contribs[idx]
		if contrib.Target == r.LastCharName {
			color = selfColor
		}
		bars = append(bars, toThroughputBar(contrib, idx, totalDamage, contribs[0].TotalDamage, durationSec, color))
	}
	for idx:=presenter.ThroughputBarCount; idx<len(contribs); idx++ {
		contrib := contribs[idx]
		if contrib.Target == r.LastCharName {
			bars = append(bars, toThroughputBar(contrib, idx, totalDamage, contribs[0].TotalDamage, durationSec, selfColor))
		}
	}
	return bars
}

func toThroughputBar(c *Contribution, index int, totalDamage int64, maxDmg int64, durationSec float64, color string) presenter.ThroughputBar {
	result := presenter.ThroughputBar{
		LeftText:   fmt.Sprintf("%d. %s", index+1, c.Target),
		CenterText: "",
		RightText: fmt.Sprintf("%s [%s]",
			parsedefs.FormatFixed(float64(c.TotalDamage)),
			parsedefs.FormatPercent(float64(c.TotalDamage)/float64(totalDamage)),
		),
		RightStyle: presenter.MonoStyle,
		Sectors:    nil,
	}
	result.Sectors=[]presenter.BarSector{{Color: color, Portion: float64(c.TotalDamage)/float64(maxDmg)}}
	return result
}

