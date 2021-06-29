package parsedefs

import (
	"bytes"
	"encoding/gob"
)

// FightReportSet represents the whole collection of reports for a fight
type FightReportSet map[string]FightReport

type FightReportFactory interface {
	// Type of report this factory creates
	Type() string

	// NewEmpty creates a report of this type focused on a fight with the specified target
	NewEmpty(target string) FightReport

	// Merge a collection of reports of this type
	Merge(reports []FightReport) FightReport

	// Deserialize a serialized report of this factory's type.
	Deserialize(serialized []byte) (FightReport, error)
}

var reportRegistry=map[string]FightReportFactory{}

func RegisterReport(factory FightReportFactory) {
	reportRegistry[factory.Type()]=factory
}

type encodedReportEntry struct {
	reportType string
	reportData []byte
}

func (rs FightReportSet) GobEncode() ([]byte, error) {
	var eres []*encodedReportEntry
	for reportType, report := range rs {
		serial, err := report.Serialize()
		if err!=nil {return nil, err}
		eres = append(eres, &encodedReportEntry{reportType, serial})
	}
	buf := &bytes.Buffer{}
	err := gob.NewEncoder(buf).Encode(eres)
	return buf.Bytes(), err
}

func (rs FightReportSet) GobDecode(encoded []byte) error {
	var eres []*encodedReportEntry
	err := gob.NewDecoder(bytes.NewReader(encoded)).Decode(&eres)
	if err!=nil {
		return err
	}
	for _, ere := range eres {
		if factory, ok := reportRegistry[ere.reportType]; ok {
			report, err := factory.Deserialize(ere.reportData)
			if err!=nil {
				return err
			} else {
				rs[ere.reportType]=report
			}
		}
	}
	return nil
}

// NewFightReports create a collection of reports specialized to a fight against the specified target
func NewFightReports(target string) FightReportSet {
	rs := FightReportSet{}
	for reportType, factory := range reportRegistry {
		rs[reportType]=factory.NewEmpty(target)
	}
	return rs
}

func MergeFightReports(sets []FightReportSet) FightReportSet {
	reportNames := map[string]struct{}{}
	for _, set := range sets {
		for repName, _ := range set {
			reportNames[repName]=struct{}{}
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
		result[repName]= reportRegistry[repName].Merge(reps)
	}
	return result
}
