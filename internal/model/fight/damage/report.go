package damage

import (
	"bytes"
	"encoding/gob"
	"github.com/gontikr99/chutzparse/internal/model/fight"
	"github.com/vugu/vugu"
)

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

func (r *Report) Serialize() ([]byte, error) {
	b := &bytes.Buffer{}
	err := gob.NewEncoder(b).Encode(r)
	return b.Bytes(), err
}

func (r *Report) Detail(fight *fight.Fight) vugu.Builder {return nil}
func (r *Report) Finalize() fight.FightReport                {return r}

type ReportFactory struct {}

func (r ReportFactory) Type() string {return "Damage"}

func (r ReportFactory) NewEmpty(target string) fight.FightReport {
	return &Report{
		Target:        target,
		Contributions: make(map[string]*Contribution),
	}
}

func (r ReportFactory) Merge(reports []fight.FightReport) fight.FightReport {
	// FIXME: implement
	return r.NewEmpty("")
}

func (r ReportFactory) Deserialize(serialized []byte) (fight.FightReport, error) {
	var result Report
	err := gob.NewDecoder(bytes.NewReader(serialized)).Decode(&result)
	return &result, err
}

