package heal

import (
	"bytes"
	"encoding/gob"
	"github.com/gontikr99/chutzparse/internal/model/fight"
	"github.com/vugu/vugu"
)

type Report struct {
	Belligerant string
	Contributions map[string]*Contribution
	LastCharName string
}

type Contribution struct {
	Source string
	TotalHealed int64
	HealByEpoch map[int]int64
}

func (r *Report) Serialize() ([]byte, error) {
	b := &bytes.Buffer{}
	err := gob.NewEncoder(b).Encode(r)
	return b.Bytes(), err
}

func (r *Report) Detail(fight *interface{}) vugu.Builder {
	// FIXME: implement
	return nil
}
func (r *Report) Finalize() fight.FightReport {return r}

type ReportFactory struct {}

func (r ReportFactory) Type() string {return "Healing"}

func (r ReportFactory) NewEmpty(target string) fight.FightReport {
	return &Report{
		Belligerant: target,
		Contributions: make(map[string]*Contribution),
	}
}

func (r ReportFactory) Merge(reports []fight.FightReport) fight.FightReport {
	// FIXME: implement
	return nil
}
func (r ReportFactory) Deserialize(serialized []byte) (fight.FightReport, error) {
	var he Report
	err := gob.NewDecoder(bytes.NewReader(serialized)).Decode(&he)
	return &he, err
}

