package heal

import (
	"github.com/gontikr99/chutzparse/internal/model/fight"
	"time"
)

type Report struct {
	Belligerent   string
	Contributions map[string]*Contribution
	LastCharName  string
	StartTime     time.Time
	EndTime       time.Time
}

type Contribution struct {
	Source      string
	TotalHealed int64
	HealByEpoch map[int]int64
}

func (r *Report) Finalize(f *fight.Fight) fight.FightReport {
	r.StartTime = f.StartTime
	r.EndTime = f.LastActivity
	return r
}

type ReportFactory struct{}

func (r ReportFactory) Type() string { return "Healing" }

func (r ReportFactory) NewEmpty(target string) fight.FightReport {
	return &Report{
		Belligerent:   target,
		Contributions: make(map[string]*Contribution),
	}
}

func (r ReportFactory) Merge(reports []fight.FightReport) fight.FightReport {
	result := &Report{
		Contributions: make(map[string]*Contribution),
	}
	if len(reports) > 0 {
		result.StartTime = reports[0].(*Report).StartTime
		result.EndTime = reports[0].(*Report).EndTime
	}
	for _, reportIf := range reports {
		report := reportIf.(*Report)
		if result.Belligerent == "" {
			result.Belligerent = report.Belligerent + " and others"
		}
		if result.LastCharName == "" {
			result.LastCharName = ""
		}
		for name, contrib := range report.Contributions {
			update, present := result.Contributions[name]
			if !present {
				update = &Contribution{Source: name, HealByEpoch: make(map[int]int64)}
				result.Contributions[name] = update
			}
			for epoch, healed := range contrib.HealByEpoch {
				resHeal, _ := update.HealByEpoch[epoch]
				if healed > resHeal {
					update.HealByEpoch[epoch] = healed
				}
			}
		}
		if report.StartTime.Before(result.StartTime) {
			result.StartTime = report.StartTime
		}
		if report.EndTime.After(result.EndTime) {
			result.EndTime = report.EndTime
		}
	}
	for _, contrib := range result.Contributions {
		contrib.TotalHealed = 0
		for _, healed := range contrib.HealByEpoch {
			contrib.TotalHealed += healed
		}
	}
	return result
}
