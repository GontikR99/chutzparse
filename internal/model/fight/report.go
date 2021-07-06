package fight

import "encoding/gob"

// FightReportSet represents the whole collection of reports for a fight
type FightReportSet map[string]FightReport

type FightReportFactory interface {
	// Type of fight this factory creates
	Type() string

	// NewEmpty creates a fight of this type focused on a fight with the specified target
	NewEmpty(target string) FightReport

	// Merge a collection of reports of this type
	Merge(reports []FightReport) FightReport
}

var reportRegistry = map[string]FightReportFactory{}

func RegisterReport(factory FightReportFactory) {
	reportRegistry[factory.Type()] = factory
	gob.RegisterName("FightReport:"+factory.Type(), factory.NewEmpty())
}

// NewFightReports create a collection of reports specialized to a fight against the specified target
func NewFightReports(target string) FightReportSet {
	rs := FightReportSet{}
	for reportType, factory := range reportRegistry {
		rs[reportType] = factory.NewEmpty(target)
	}
	return rs
}

func MergeFightReports(sets []FightReportSet) FightReportSet {
	reportNames := map[string]struct{}{}
	for _, set := range sets {
		for repName, _ := range set {
			reportNames[repName] = struct{}{}
		}
	}
	result := FightReportSet{}
	for repName, _ := range reportNames {
		var reps []FightReport
		for _, set := range sets {
			if report, ok := set[repName]; ok {
				reps = append(reps, report)
			}
		}
		result[repName] = reportRegistry[repName].Merge(reps)
	}
	return result
}
