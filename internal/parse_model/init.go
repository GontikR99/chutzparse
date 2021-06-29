package parse_model

import (
	"github.com/gontikr99/chutzparse/internal/parse_model/damage"
	"github.com/gontikr99/chutzparse/internal/parse_model/parsedefs"
)

// Start starts th
func init() {
	parsedefs.RegisterReport(damage.ReportFactory{})
}