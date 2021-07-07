// +build wasm,electron

package heal

import (
	"fmt"
	"github.com/gontikr99/chutzparse/internal/model/fight"
	"github.com/gontikr99/chutzparse/internal/model/iff"
	"github.com/gontikr99/chutzparse/internal/model/parsedefs"
	"github.com/gontikr99/chutzparse/internal/presenter"
	"sort"
	"strings"
	"time"
)

type contribByHealRev []*Contribution

func (c contribByHealRev) Len() int           { return len(c) }
func (c contribByHealRev) Less(i, j int) bool { return c[i].TotalHealed > c[j].TotalHealed }
func (c contribByHealRev) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }

type characterAndWards struct {
	Attribution  string
	TotalHealed  int64
	Contributors []*Contribution
}

func (c *characterAndWards) toThroughputBar(index int, totalHealed int64, maxHealed int64, durationSec float64, colors []string) presenter.ThroughputBar {
	dps := float64(c.TotalHealed) / durationSec
	result := presenter.ThroughputBar{
		LeftText:   fmt.Sprintf("%d. %s", index+1, c.Attribution),
		CenterText: "",
		RightText: fmt.Sprintf("%s [%.3g%%] = %s hps",
			parsedefs.RenderAmount(float64(c.TotalHealed)),
			float64(100*c.TotalHealed)/float64(totalHealed),
			parsedefs.RenderAmount(dps)),
		Sectors: nil,
	}
	if len(c.Contributors) > 1 {
		result.LeftText = result.LeftText + " + wards"
	}
	for idx, ctb := range c.Contributors {
		result.Sectors = append(result.Sectors, presenter.BarSector{
			Color:   colors[idx%len(colors)],
			Portion: float64(ctb.TotalHealed) / float64(maxHealed),
		})
	}
	return result
}

type cnpByHealedRev []*characterAndWards

func (c cnpByHealedRev) Len() int           { return len(c) }
func (c cnpByHealedRev) Less(i, j int) bool { return c[i].TotalHealed > c[j].TotalHealed }
func (c cnpByHealedRev) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }

const throughputBarCount = 10

func (r *Report) Throughput(fight *fight.Fight) []presenter.ThroughputBar {
	duration := fight.LastActivity.Sub(fight.StartTime)
	if duration <= 0 {
		duration = 1 * time.Second
	}

	var totalHealed int64
	for _, contrib := range r.Contributions {
		totalHealed += contrib.TotalHealed
	}

	durationSec := float64(duration) / float64(time.Second)
	hps := float64(totalHealed) / durationSec

	var bars []presenter.ThroughputBar
	bars = append(bars, presenter.ThroughputBar{
		LeftText:  fmt.Sprintf("[Healing] %s in %s", r.Belligerent, parsedefs.RenderAmount(durationSec)),
		RightText: fmt.Sprintf("%s = %s hps", parsedefs.RenderAmount(float64(totalHealed)), parsedefs.RenderAmount(hps)),
		Sectors:   []presenter.BarSector{{"DimGray", 1.0}},
	})

	if len(r.Contributions) == 0 {
		return bars
	}
	var cnps []*characterAndWards
	for _, contributor := range r.Contributions {
		attr := contributor.Source
		if owner := iff.GetOwner(contributor.Source); owner != "" && strings.HasSuffix(contributor.Source, "`s ward") {
			attr = owner
		}
		var update *characterAndWards
		for i := range cnps {
			if cnps[i].Attribution == attr {
				update = cnps[i]
			}
		}
		if update == nil {
			update = &characterAndWards{
				Attribution:  attr,
				TotalHealed:  0,
				Contributors: nil,
			}
			cnps = append(cnps, update)
		}
		update.TotalHealed += contributor.TotalHealed
		update.Contributors = append(update.Contributors, contributor)
	}

	sort.Sort(cnpByHealedRev(cnps))
	for _, cnp := range cnps {
		sort.Sort(contribByHealRev(cnp.Contributors))

	}

	maxDmg := cnps[0].TotalHealed
	for i := 0; i < throughputBarCount && i < len(cnps); i++ {
		colors := []string{"blue", "lightblue"}
		if cnps[i].Attribution == r.LastCharName {
			colors = []string{"limegreen", "forestgreen"}
		}
		bars = append(bars, cnps[i].toThroughputBar(i, totalHealed, maxDmg, durationSec, colors))
	}
	// Add me if I haven't already been shown
	for i := throughputBarCount; i < len(cnps); i++ {
		if cnps[i].Attribution == r.LastCharName {
			colors := []string{"limegreen", "forestgreen"}
			bars = append(bars, cnps[i].toThroughputBar(i, totalHealed, maxDmg, durationSec, colors))
		}
	}
	return bars
}
