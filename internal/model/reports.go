package model

import (
	"github.com/gontikr99/chutzparse/internal/model/fight"
	"github.com/gontikr99/chutzparse/internal/model/fight/damage"
	"github.com/gontikr99/chutzparse/internal/model/fight/heal"
)

func RegisterReports() {
	fight.RegisterReport(damage.ReportFactory{})
	fight.RegisterReport(heal.ReportFactory{})
}
