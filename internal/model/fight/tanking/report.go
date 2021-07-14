package tanking

import (
	"github.com/gontikr99/chutzparse/internal/model/fight"
	"github.com/gontikr99/chutzparse/pkg/algorithm"
	"sort"
)

type Report struct {
	Source string
	LastCharName string
	Contributions map[string]*Contribution
	ActivitySet algorithm.TimeIntervalSet
}

func NewReport(source string) *Report {
	return &Report{
		Source:        source,
		Contributions: map[string]*Contribution{},
	}
}

func (r *Report) Finalize(f *fight.Fight) fight.FightReport {
	r.ActivitySet = algorithm.NewTimeInterval(f.StartTime, f.LastActivity)
	return r
}

func (r *Report) TotalDamage() int64 {
	var totalDamage int64
	for _, contrib := range r.Contributions {
		totalDamage+=contrib.TotalDamage
	}
	return totalDamage
}

func (r *Report) Interesting() bool {
	return r.TotalDamage()!=0
}

func (r *Report) Participants(p map[string]struct{}) {
	for k := range r.Contributions {
		p[k]=struct{}{}
	}
}

func (r *Report) ContributionOf(target string) *Contribution {
	update, ok := r.Contributions[target]
	if !ok {
		update=NewContribution(target)
		r.Contributions[target]=update
	}
	return update
}

func (r *Report) SortedContributors() []*Contribution {
	var contribs []*Contribution
	for _, contrib := range r.Contributions {
		contribs = append(contribs, contrib)
	}
	sort.Sort(byTotalDamageRev(contribs))
	return contribs
}
type byTotalDamageRev []*Contribution
func (b byTotalDamageRev) Len() int {return len(b)}
func (b byTotalDamageRev) Less(i, j int) bool {return b[i].TotalDamage > b[j].TotalDamage}
func (b byTotalDamageRev) Swap(i, j int) {b[i], b[j] = b[j], b[i]}


type Contribution struct {
	Target string
	TotalDamage int64
	Hits int
}

func NewContribution(target string) *Contribution {
	return &Contribution{
		Target:      target,
	}
}

type ReportFactory struct {}
func (rf ReportFactory) Type() string {return "Tanking"}
func (rf ReportFactory) NewEmpty(source string) fight.FightReport {return NewReport(source)}
func (rf ReportFactory) Merge(reportsIf []fight.FightReport) fight.FightReport {
	if len(reportsIf)==0 {return nil}
	if len(reportsIf)==1 {return reportsIf[0]}

	result := NewReport(reportsIf[0].(*Report).Source+" and others")
	for _, reportIf := range reportsIf {
		report := reportIf.(*Report)
		if result.LastCharName=="" {
			result.LastCharName=report.LastCharName
		}
		result.ActivitySet = algorithm.UnionTimeIntervalSets(result.ActivitySet, report.ActivitySet)
		for name, contrib := range report.Contributions {
			resContrib := result.ContributionOf(name)
			resContrib.Hits += contrib.Hits
			resContrib.TotalDamage += contrib.TotalDamage
		}
	}
	return result
}