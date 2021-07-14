package model

import (
	"github.com/gontikr99/chutzparse/internal/model/fight"
	"github.com/gontikr99/chutzparse/internal/model/fight/damage"
	"github.com/gontikr99/chutzparse/internal/model/fight/heal"
	"github.com/gontikr99/chutzparse/internal/model/fight/tanking"
)

func RegisterReports() {
	fight.RegisterReport(damage.ReportFactory{})
	fight.RegisterReport(heal.ReportFactory{})
	fight.RegisterReport(tanking.ReportFactory{})
}
