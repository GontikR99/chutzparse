//go:build wasm && electron
// +build wasm,electron

package heal

import (
	"fmt"
	"github.com/gontikr99/chutzparse/internal/model/fight"
	"github.com/gontikr99/chutzparse/internal/model/parsedefs"
	"github.com/gontikr99/chutzparse/internal/presenter"
	"time"
)

func toThroughputBar(c *aggregateContributor, index int, totalHealed int64, maxHealed int64, durationSec float64, colors []string) presenter.ThroughputBar {
	dps := float64(c.TotalHealed) / durationSec
	result := presenter.ThroughputBar{
		LeftText:   fmt.Sprintf("%d. %s", index+1, c.DisplayName()),
		CenterText: "",
		RightText: fmt.Sprintf("%s [%s] = %s hps",
			parsedefs.FormatFixed(float64(c.TotalHealed)),
			parsedefs.FormatPercent(float64(c.TotalHealed)/float64(totalHealed)),
			parsedefs.FormatFixed(dps)),
		RightStyle: presenter.MonoStyle,
		Sectors:    nil,
	}
	for idx, ctb := range c.RawContributions {
		result.Sectors = append(result.Sectors, presenter.BarSector{
			Color:   colors[idx%len(colors)],
			Portion: float64(ctb.TotalHealed) / float64(maxHealed),
		})
	}
	return result
}

var otherColors = []string{"blue", "lightblue"}
var selfColors = []string{"limegreen", "forestgreen"}

func (r *Report) Throughput(fight *fight.Fight) []presenter.ThroughputBar {
	duration := fight.LastActivity.Sub(fight.StartTime)
	if duration <= 1 {
		duration = 1 * time.Second
	}
	durationSec := float64(duration) / float64(time.Second)

	aggRep := r.Aggregate()
	hps := float64(aggRep.TotalHealed) / durationSec

	var bars []presenter.ThroughputBar
	bars = append(bars, presenter.ThroughputBar{
		LeftText:   fmt.Sprintf("[Healing] %s in %ss", r.Belligerent, parsedefs.FormatFixed(durationSec)),
		RightText:  fmt.Sprintf("%s = %s hps", parsedefs.FormatFixed(float64(aggRep.TotalHealed)), parsedefs.FormatFixed(hps)),
		RightStyle: presenter.MonoStyle,
		Sectors:    []presenter.BarSector{{"dimgray", 1.0}},
	})

	if len(aggRep.Contributions) == 0 {
		return bars
	}
	agContribs := aggRep.SortedContributors()
	maxHealed := agContribs[0].TotalHealed

	for i := 0; i < presenter.ThroughputBarCount && i < len(agContribs); i++ {
		colors := otherColors
		if agContribs[i].AttributedSource == r.LastCharName {
			colors = selfColors
		}
		bars = append(bars, toThroughputBar(agContribs[i], i, aggRep.TotalHealed, maxHealed, durationSec, colors))
	}
	for i := presenter.ThroughputBarCount; i < len(agContribs); i++ {
		if agContribs[i].AttributedSource == r.LastCharName {
			bars = append(bars, toThroughputBar(agContribs[i], i, aggRep.TotalHealed, maxHealed, durationSec, selfColors))
		}
	}
	return bars
}
