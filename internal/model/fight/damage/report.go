package damage

import (
	"github.com/gontikr99/chutzparse/internal/model/boosts"
	"github.com/gontikr99/chutzparse/internal/model/fight"
	"github.com/gontikr99/chutzparse/pkg/algorithm"
	"sort"
	"strings"
)

// Report is the damage dealt/HPS report
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
	return len(r.Contributions) != 0
}

func (r *Report) Participants(p map[string]struct{}) {
	for k, contrib := range r.Contributions {
		if contrib.TotalDamage != 0 {
			p[k] = struct{}{}
		}
	}
}

// NewReport creates a new empty report with the specified target
func NewReport(target string) *Report {
	return &Report{Target: target, Contributions: make(map[string]*Contribution), ActivitySet: algorithm.EmptyTimeIntervalSet}
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
	Deaths      map[int]struct{}
	Boosts      boosts.BoostSet
}

func NewContribution(source string) *Contribution {
	return &Contribution{
		Source:      source,
		Categorized: map[string]*Category{},
		Deaths:      map[int]struct{}{},
		Boosts:      boosts.BoostSet{},
	}
}

type flagEntry struct {
	logId    int
	flagText string
}
type byLogId []flagEntry

func (b byLogId) Len() int           { return len(b) }
func (b byLogId) Less(i, j int) bool { return b[i].logId < b[j].logId }
func (b byLogId) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }

func (c *Contribution) Flags() string {
	var flagEntries []flagEntry
	for logId := range c.Deaths {
		flagEntries = append(flagEntries, flagEntry{logId, "X"})
	}
	for boostType, idMap := range c.Boosts {
		for logId := range idMap {
			flagEntries = append(flagEntries, flagEntry{logId, boostType.String()})
		}
	}
	sort.Sort(byLogId(flagEntries))
	sb := &strings.Builder{}
	for _, entry := range flagEntries {
		sb.WriteString(entry.flagText)
	}
	return sb.String()
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
	for _, reportIf := range reports {
		report := reportIf.(*Report)
		if result.LastCharName == "" {
			result.LastCharName = report.LastCharName
		}
		result.ActivitySet = algorithm.UnionTimeIntervalSets(result.ActivitySet, report.ActivitySet)
		for name, contrib := range report.Contributions {
			resContrib := result.ContributionOf(name)
			resContrib.Boosts.AddAll(contrib.Boosts)
			resContrib.TotalDamage += contrib.TotalDamage
			for id := range contrib.Deaths {
				resContrib.Deaths[id] = struct{}{}
			}
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
