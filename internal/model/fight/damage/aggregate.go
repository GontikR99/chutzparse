package damage

import (
	iff2 "github.com/gontikr99/chutzparse/internal/iff"
	"sort"
)

// Aggregate combines the contributions of pets with those of their owners
func (r *Report) Aggregate() *aggregateReport {
	agRep := newAggregateReport(r.Target)
	for source, contrib := range r.Contributions {
		agContrib := agRep.ContributionOf(source)
		agContrib.TotalDamage += contrib.TotalDamage
		agRep.TotalDamage += contrib.TotalDamage
		agContrib.RawContributions = append(agContrib.RawContributions, contrib)
		for dispName, cat := range contrib.Categorized {
			agCat := agContrib.CategoryOf(source, dispName)
			agCat.TotalDamage += cat.TotalDamage
			agCat.Success += cat.Success
			agCat.Failure += cat.Failure
		}
	}
	for _, agContrib := range agRep.Contributions {
		sort.Sort(contribByDamageRev(agContrib.RawContributions))
	}
	return agRep
}

type aggregateReport struct {
	Target        string
	TotalDamage   int64
	Contributions map[string]*aggregateContributor
}

func newAggregateReport(target string) *aggregateReport {
	return &aggregateReport{Target: target, Contributions: map[string]*aggregateContributor{}}
}

func (ar *aggregateReport) ContributionOf(source string) *aggregateContributor {
	attributedSource := source
	if owner := iff2.GetOwner(source); owner != "" {
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
	sort.Sort(acByDamageRev(contribs))
	return contribs
}

type aggregateContributor struct {
	AttributedSource string
	Sources          map[string]struct{}
	TotalDamage      int64
	Categorized      map[string]*Category
	RawContributions []*Contribution
}

func newAggregateContributor(attributedSource string, source string) *aggregateContributor {
	return &aggregateContributor{
		AttributedSource: attributedSource,
		Sources:          map[string]struct{}{source: {}},
		Categorized:      map[string]*Category{},
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
		return ac.AttributedSource + " + pets"
	}
	return ac.AttributedSource
}

func (ac *aggregateContributor) CategoryOf(source string, displayName string) *Category {
	if owner := iff2.GetOwner(source); owner != "" {
		displayName = displayName + " (" + source + ")"
	}
	update, ok := ac.Categorized[displayName]
	if !ok {
		update = NewCategory(displayName)
		ac.Categorized[displayName] = update
	}
	return update
}

func (ac *aggregateContributor) SortedCategories() []*Category {
	var result []*Category
	for _, cat := range ac.Categorized {
		result = append(result, cat)
	}
	sort.Sort(catByDamageRev(result))
	return result
}

type acByDamageRev []*aggregateContributor

func (a acByDamageRev) Len() int           { return len(a) }
func (a acByDamageRev) Less(i, j int) bool { return a[i].TotalDamage > a[j].TotalDamage }
func (a acByDamageRev) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
