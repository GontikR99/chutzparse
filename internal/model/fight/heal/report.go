package heal

import (
	"github.com/gontikr99/chutzparse/internal/model/fight"
	"github.com/gontikr99/chutzparse/pkg/algorithm"
	"github.com/gontikr99/chutzparse/pkg/console"
)

type Report struct {
	Belligerent   string
	Contributions map[string]*Contribution
	LastCharName  string
	ActivitySet   algorithm.TimeIntervalSet
}

func (r *Report) Finalize(f *fight.Fight) fight.FightReport {
	r.ActivitySet = algorithm.NewTimeInterval(f.StartTime, f.LastActivity)
	return r
}

func (r *Report) Interesting() bool {
	return false
}

func NewReport(belligerent string) *Report {
	return &Report{Belligerent: belligerent, Contributions: map[string]*Contribution{}}
}

// ContributionOf returns a pointer to the contribution record for a specified source, adding such a record to the
// the report if none yet exists
func (r *Report) ContributionOf(source string) *Contribution {
	update, ok := r.Contributions[source]
	if !ok {
		update = NewContribution(source)
		r.Contributions[source] = update
	}
	return update
}

type Contribution struct {
	Source      string
	TotalHealed int64
	Categorized map[string]*Category
}

func NewContribution(source string) *Contribution {
	return &Contribution{Source: source, Categorized: make(map[string]*Category)}
}

func (c *Contribution) CategoryOf(displayName string) *Category {
	update, ok := c.Categorized[displayName]
	if !ok {
		update = NewCategory(displayName)
		c.Categorized[displayName] = update
	}
	return update
}

type contribByHealingRev []*Contribution

func (c contribByHealingRev) Len() int           { return len(c) }
func (c contribByHealingRev) Less(i, j int) bool { return c[i].TotalHealed > c[j].TotalHealed }
func (c contribByHealingRev) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }

type Category struct {
	DisplayName   string
	HealedByEpoch map[int]*EpochStat
}

func NewCategory(displayName string) *Category {
	return &Category{DisplayName: displayName, HealedByEpoch: map[int]*EpochStat{}}
}

func (c *Category) EpochOf(epoch int) *EpochStat {
	update, ok := c.HealedByEpoch[epoch]
	if !ok {
		update = NewEpochStat()
		c.HealedByEpoch[epoch] = update
	}
	return update
}

func (c *Category) TotalStats() (count int, healed int64) {
	for _, es := range c.HealedByEpoch {
		count += es.Count
		healed += es.TotalHealed
	}
	return
}

type EpochStat struct {
	Count       int
	TotalHealed int64
}

func NewEpochStat() *EpochStat {
	return &EpochStat{}
}

type ReportFactory struct{}

func (r ReportFactory) Type() string { return "Healing" }

func (r ReportFactory) NewEmpty(target string) fight.FightReport {
	return NewReport(target)
}

func (r ReportFactory) Merge(reports []fight.FightReport) fight.FightReport {
	if len(reports) == 0 {
		return nil
	}
	if len(reports) == 1 {
		return reports[0]
	}

	result := NewReport(reports[0].(*Report).Belligerent + " and others")
	for _, reportIf := range reports {
		report := reportIf.(*Report)
		if result.LastCharName == "" {
			result.LastCharName = report.LastCharName
		}
		result.ActivitySet = algorithm.UnionTimeIntervalSets(result.ActivitySet, report.ActivitySet)
		for name, contrib := range report.Contributions {
			resContrib := result.ContributionOf(name)
			for catName, cat := range contrib.Categorized {
				resCat := resContrib.CategoryOf(catName)
				for epoch, epochStat := range cat.HealedByEpoch {
					resEpochStat := resCat.EpochOf(epoch)
					if epochStat.TotalHealed > resEpochStat.TotalHealed {
						resEpochStat.TotalHealed = epochStat.TotalHealed
					}
					if epochStat.Count > resEpochStat.Count {
						resEpochStat.Count = epochStat.Count
					}
				}
			}
		}
	}
	for _, resContrib := range result.Contributions {
		for _, resCat := range resContrib.Categorized {
			for _, resEpochStat := range resCat.HealedByEpoch {
				resContrib.TotalHealed += resEpochStat.TotalHealed
			}
		}
	}
	console.Log("merged: ", result.ActivitySet)
	return result
}
