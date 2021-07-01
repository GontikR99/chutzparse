package model

import (
	"github.com/gontikr99/chutzparse/internal/model/fight"
	"github.com/gontikr99/chutzparse/internal/model/fight/damage"
)

// Start starts th
func init() {
	fight.RegisterReport(damage.ReportFactory{})
}