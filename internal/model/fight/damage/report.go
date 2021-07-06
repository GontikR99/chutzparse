package damage

import (
	"github.com/gontikr99/chutzparse/internal/model/fight"
)

type Report struct {
	Target        string
	LastCharName  string
	Contributions map[string]*Contribution
}

type Contribution struct {
	Source      string
	TotalDamage int64
}

func (c *Contribution) DamageTotal() int64 {
	return c.TotalDamage
}

func (r *Report) Finalize() fight.FightReport { return r }

type ReportFactory struct{}

func (r ReportFactory) Type() string { return "Damage" }

func (r ReportFactory) NewEmpty(target string) fight.FightReport {
	return &Report{
		Target:        target,
		Contributions: make(map[string]*Contribution),
	}
}

func (r ReportFactory) Merge(reports []fight.FightReport) fight.FightReport {
	result := &Report{
		Contributions: make(map[string]*Contribution),
	}
	if len(reports) == 0 {
		return result
	}
	result.Target = reports[0].(*Report).Target + " and others"
	for _, reportIf := range reports {
		report := reportIf.(*Report)
		if result.LastCharName == "" {
			result.LastCharName = report.LastCharName
		}
		for name, contrib := range report.Contributions {
			update, present := result.Contributions[name]
			if !present {
				update = &Contribution{Source: name}
				result.Contributions[name] = update
			}
			update.TotalDamage += contrib.TotalDamage
		}
	}
	return result
}
