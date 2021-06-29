package damage

import (
	"bytes"
	"encoding/gob"
	"github.com/gontikr99/chutzparse/internal/eqlog"
	"github.com/gontikr99/chutzparse/internal/parse_model/parsedefs"
	"github.com/vugu/vugu"
)

const throughputBarCount=10

type Report struct {
	Target string
	LastCharName string
	Contributions map[string]*Contribution
}

type Contribution struct {
	Source string
	TotalDamage int64
}

func (c *Contribution) DamageTotal() int64 {
	return c.TotalDamage
}

func (r *Report) Offer(entry *eqlog.LogEntry, epoch int) parsedefs.FightReport {
	r.LastCharName = entry.Character
	dmg, ok := entry.Meaning.(*eqlog.DamageLog)
	if !ok {return r}
	if dmg.Target!=r.Target && dmg.Target!=r.Target+"`s pet" && dmg.Target!=r.Target+"`s warder" {
		return r
	}
	if _, ok := r.Contributions[dmg.Source]; !ok {
		r.Contributions[dmg.Source]=&Contribution{Source: dmg.Source}
	}
	r.Contributions[dmg.Source].TotalDamage+=dmg.Amount
	return r
}

func (r *Report) Serialize() ([]byte, error) {
	b := &bytes.Buffer{}
	err := gob.NewEncoder(b).Encode(r)
	return b.Bytes(), err
}

func (r *Report) Detail(fight *parsedefs.Fight) vugu.Builder {return nil}
func (r *Report) Finalize() parsedefs.FightReport            {return r}

type ReportFactory struct {}

func (r ReportFactory) Type() string {return "Damage"}

func (r ReportFactory) NewEmpty(target string) parsedefs.FightReport {
	return &Report{
		Target:        target,
		Contributions: make(map[string]*Contribution),
	}
}

func (r ReportFactory) Merge(reports []parsedefs.FightReport) parsedefs.FightReport {
	// FIXME: implement
	return r.NewEmpty("")
}

func (r ReportFactory) Deserialize(serialized []byte) (parsedefs.FightReport, error) {
	var result Report
	err := gob.NewDecoder(bytes.NewReader(serialized)).Decode(&result)
	return &result, err
}

