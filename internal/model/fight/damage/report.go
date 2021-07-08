package damage

import (
	"github.com/gontikr99/chutzparse/internal/model/fight"
	"github.com/gontikr99/chutzparse/pkg/algorithm"
)

// Report is the damage dealt/DPS report
type Report struct {
	Target        string
	LastCharName  string
	Contributions map[string]*Contribution
	ActivitySet   algorithm.TimeIntervalSet
}

func (r *Report) Finalize(f *fight.Fight) fight.FightReport {
	r.ActivitySet = algorithm.NewTimeInterval(f.StartTime, f.LastActivity)
	return r
}

func (r *Report) Interesting() bool {
	return len(r.Contributions)!=0
}

// NewReport creates a new empty report with the specified target
func NewReport(target string) *Report {
	return &Report{Target: target, Contributions: make(map[string]*Contribution)}
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
	TotalDamage int64
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

type contribByDamageRev []*Contribution

func (c contribByDamageRev) Len() int           { return len(c) }
func (c contribByDamageRev) Less(i, j int) bool { return c[i].TotalDamage > c[j].TotalDamage }
func (c contribByDamageRev) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }

type Category struct {
	DisplayName string
	Success     int
	Failure     int
	TotalDamage int64
}

func NewCategory(displayName string) *Category {
	return &Category{DisplayName: displayName}
}

type catByDamageRev []*Category

func (c catByDamageRev) Len() int           { return len(c) }
func (c catByDamageRev) Less(i, j int) bool { return c[i].TotalDamage > c[j].TotalDamage }
func (c catByDamageRev) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }

type ReportFactory struct{}

func (r ReportFactory) Type() string { return "Damage" }

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

	result := NewReport(reports[0].(*Report).Target + " and others")
	result.ActivitySet = reports[0].(*Report).ActivitySet
	for _, reportIf := range reports {
		report := reportIf.(*Report)
		if result.LastCharName == "" {
			result.LastCharName = report.LastCharName
		}
		result.ActivitySet = algorithm.UnionTimeIntervalSets(result.ActivitySet, report.ActivitySet)
		for name, contrib := range report.Contributions {
			resContrib := result.ContributionOf(name)
			resContrib.TotalDamage += contrib.TotalDamage
			for catName, cat := range contrib.Categorized {
				resCat := resContrib.CategoryOf(catName)
				resCat.TotalDamage += cat.TotalDamage
				resCat.Success += cat.Success
				resCat.Failure += cat.Failure
			}
		}
	}
	return result
}
