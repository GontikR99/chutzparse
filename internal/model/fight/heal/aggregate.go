package heal

import (
	"github.com/gontikr99/chutzparse/internal/model/iff"
	"sort"
)

// Aggregate combines the contributions of wards with those of their owners
func (r *Report) Aggregate() *aggregateReport {
	agRep := newAggregateReport(r.Belligerent)
	for source, contrib := range r.Contributions {
		agContrib := agRep.ContributionOf(source)
		agContrib.TotalHealed += contrib.TotalHealed
		agRep.TotalHealed += contrib.TotalHealed
		agContrib.RawContributions = append(agContrib.RawContributions, contrib)
		for dispName, cat := range contrib.Categorized {
			agCat := agContrib.CategoryOf(source, dispName)
			count, healed := cat.TotalStats()
			agCat.Count += count
			agCat.TotalHealed += healed
		}
	}
	for _, agContrib := range agRep.Contributions {
		sort.Sort(contribByHealingRev(agContrib.RawContributions))
	}
	return agRep
}

type aggregateReport struct {
	Belligerent   string
	TotalHealed   int64
	Contributions map[string]*aggregateContributor
}

func newAggregateReport(belligerent string) *aggregateReport {
	return &aggregateReport{
		Belligerent:   belligerent,
		TotalHealed:   0,
		Contributions: map[string]*aggregateContributor{},
	}
}

func (ar *aggregateReport) ContributionOf(source string) *aggregateContributor {
	attributedSource := source
	if owner := iff.GetOwner(source); owner != "" {
		attributedSource = owner
	}
	update, ok := ar.Contributions[attributedSource]
	if !ok {
		update = newAggregateContributor(attributedSource, source)
		ar.Contributions[attributedSource] = update
	}
	update.Sources[source] = struct{}{}
	return update
}

// SortedContributors returns an ordered list of the contributors, sorted from highest total damage to lowest.
func (ar *aggregateReport) SortedContributors() []*aggregateContributor {
	var contribs []*aggregateContributor
	for _, contrib := range ar.Contributions {
		contribs = append(contribs, contrib)
	}
	sort.Sort(acByHealingRev(contribs))
	return contribs
}

type aggregateContributor struct {
	AttributedSource string
	Sources          map[string]struct{}
	TotalHealed      int64
	Categorized      map[string]*aggregateCategory
	RawContributions []*Contribution
}

func newAggregateContributor(attributedSource string, source string) *aggregateContributor {
	return &aggregateContributor{
		AttributedSource: attributedSource,
		Sources:          map[string]struct{}{source: {}},
		Categorized:      map[string]*aggregateCategory{},
	}
}

// DisplayName returns the name we should display for this contributor.  Usually that's
// the character's name, but if only one of the character's pet has been detected, show the pet instead.
func (ac *aggregateContributor) DisplayName() string {
	if len(ac.Sources) == 1 {
		for name := range ac.Sources {
			return name
		}
	} else {
		return ac.AttributedSource + " + wards"
	}
	return ac.AttributedSource
}

func (ac *aggregateContributor) CategoryOf(source string, displayName string) *aggregateCategory {
	if owner := iff.GetOwner(source); owner != "" {
		displayName = displayName + " (" + source + ")"
	}
	update, ok := ac.Categorized[displayName]
	if !ok {
		update = newAggregateCategory(displayName)
		ac.Categorized[displayName] = update
	}
	return update
}

func (ac *aggregateContributor) SortedCategories() []*aggregateCategory {
	var result []*aggregateCategory
	for _, cat := range ac.Categorized {
		result = append(result, cat)
	}
	sort.Sort(acatByHealingRev(result))
	return result
}

type acByHealingRev []*aggregateContributor

func (a acByHealingRev) Len() int           { return len(a) }
func (a acByHealingRev) Less(i, j int) bool { return a[i].TotalHealed > a[j].TotalHealed }
func (a acByHealingRev) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type aggregateCategory struct {
	DisplayName string
	TotalHealed int64
	Count       int
}

func newAggregateCategory(displayName string) *aggregateCategory {
	return &aggregateCategory{DisplayName: displayName}
}

type acatByHealingRev []*aggregateCategory

func (a acatByHealingRev) Len() int           { return len(a) }
func (a acatByHealingRev) Less(i, j int) bool { return a[i].TotalHealed > a[j].TotalHealed }
func (a acatByHealingRev) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
