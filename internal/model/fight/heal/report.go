package heal

import (
	"github.com/gontikr99/chutzparse/internal/model/fight"
)

type Report struct {
	Belligerent   string
	Contributions map[string]*Contribution
	LastCharName  string
}

type Contribution struct {
	Source      string
	TotalHealed int64
	HealByEpoch map[int]int64
}

func (r *Report) Finalize() fight.FightReport { return r }

type ReportFactory struct{}

func (r ReportFactory) Type() string { return "Healing" }

func (r ReportFactory) NewEmpty(target string) fight.FightReport {
	return &Report{
		Belligerent:   target,
		Contributions: make(map[string]*Contribution),
	}
}

func (r ReportFactory) Merge(reports []fight.FightReport) fight.FightReport {
	// FIXME: implement
	return nil
}
