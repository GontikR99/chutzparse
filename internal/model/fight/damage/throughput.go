package damage

import (
	"fmt"
	"github.com/gontikr99/chutzparse/internal/model/fight"
	"github.com/gontikr99/chutzparse/internal/model/iff"
	"github.com/gontikr99/chutzparse/internal/model/parsedefs"
	"github.com/gontikr99/chutzparse/internal/presenter"
	"sort"
	"time"
)

type contribByDamageRev []*Contribution

func (c contribByDamageRev) Len() int           { return len(c) }
func (c contribByDamageRev) Less(i, j int) bool { return c[i].DamageTotal() > c[j].DamageTotal() }
func (c contribByDamageRev) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }

type characterAndPets struct {
	Attribution  string
	TotalDamage  int64
	Contributors []*Contribution
}

func (c *characterAndPets) toThroughputBar(index int, totalDamage int64, maxDmg int64, durationSec float64, colors []string) presenter.ThroughputBar {
	dps := float64(c.TotalDamage) / durationSec
	result := presenter.ThroughputBar{
		LeftText:   fmt.Sprintf("%d. %s", index+1, c.Attribution),
		CenterText: "",
		RightText: fmt.Sprintf("%s [%.3g%%] = %s dps",
			parsedefs.RenderAmount(float64(c.TotalDamage)),
			float64(100*c.TotalDamage)/float64(totalDamage),
			parsedefs.RenderAmount(dps)),
		Sectors: nil,
	}
	if len(c.Contributors) > 1 {
		result.LeftText = result.LeftText + " + pets"
	}
	for idx, ctb := range c.Contributors {
		result.Sectors = append(result.Sectors, presenter.BarSector{
			Color:   colors[idx%len(colors)],
			Portion: float64(ctb.DamageTotal()) / float64(maxDmg),
		})
	}
	return result
}

type cnpByDamageRev []*characterAndPets

func (c cnpByDamageRev) Len() int           { return len(c) }
func (c cnpByDamageRev) Less(i, j int) bool { return c[i].TotalDamage > c[j].TotalDamage }
func (c cnpByDamageRev) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }

const throughputBarCount = 10

func (r *Report) Throughput(fight *fight.Fight) []presenter.ThroughputBar {
	duration := fight.LastActivity.Sub(fight.StartTime)
	if duration <= 0 {
		duration = 1 * time.Second
	}

	var totalDamage int64
	for _, contrib := range r.Contributions {
		totalDamage += contrib.DamageTotal()
	}

	durationSec := float64(duration) / float64(time.Second)
	dps := float64(totalDamage) / durationSec

	var bars []presenter.ThroughputBar
	bars = append(bars, presenter.ThroughputBar{
		LeftText:  fmt.Sprintf("[Damage] %s in %ss", r.Target, parsedefs.RenderAmount(durationSec)),
		RightText: fmt.Sprintf("%s = %s dps", parsedefs.RenderAmount(float64(totalDamage)), parsedefs.RenderAmount(dps)),
		Sectors:   []presenter.BarSector{{"dimgray", 1.0}},
	})

	if len(r.Contributions) == 0 {
		return bars
	}
	var cnps []*characterAndPets
	for _, contributor := range r.Contributions {
		attr := contributor.Source
		if owner := iff.GetOwner(contributor.Source); owner != "" {
			attr = owner
		}
		var update *characterAndPets
		for i, _ := range cnps {
			if cnps[i].Attribution == attr {
				update = cnps[i]
			}
		}
		if update == nil {
			update = &characterAndPets{
				Attribution:  attr,
				TotalDamage:  0,
				Contributors: nil,
			}
			cnps = append(cnps, update)
		}
		update.TotalDamage += contributor.DamageTotal()
		update.Contributors = append(update.Contributors, contributor)
	}

	sort.Sort(cnpByDamageRev(cnps))
	for _, cnp := range cnps {
		sort.Sort(contribByDamageRev(cnp.Contributors))

	}

	maxDmg := cnps[0].TotalDamage
	for i := 0; i < throughputBarCount && i < len(cnps); i++ {
		colors := []string{"blue", "lightblue"}
		if cnps[i].Attribution == r.LastCharName {
			colors = []string{"red", "darksalmon"}
		}
		bars = append(bars, cnps[i].toThroughputBar(i, totalDamage, maxDmg, durationSec, colors))
	}
	// Add me if I haven't already been shown
	for i := throughputBarCount; i < len(cnps); i++ {
		if cnps[i].Attribution == r.LastCharName {
			colors := []string{"red", "darksalmon"}
			bars = append(bars, cnps[i].toThroughputBar(i, totalDamage, maxDmg, durationSec, colors))
		}
	}
	return bars
}
