package dmodel

import "time"

type HealPerf struct {
	FirstHeal time.Time
	LastHeal time.Time
	Total int64
	Actual int64
}

