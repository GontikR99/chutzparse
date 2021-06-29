package parsedefs

import (
	"time"
)

// HitEvent represents a single strike or heal that is shown on the hit display
type HitEvent struct {
	Text  string
	Color string
	Big   bool
}

const ChannelHitTop = "hitDisplayTop"
const ChannelHitBottom = "hitDisplayBottom"

type ThroughputState struct {
	FightId int
	TopBars []ThroughputBar
	BottomBars []ThroughputBar
}

type ThroughputStateEvent struct {
	Content []ThroughputState
}

const ChannelThroughput = "throughputDisplay"

// ThroughputBar is a single bar shown on the throughput display.
type ThroughputBar struct {
	LeftText   string
	CenterText string
	RightText  string
	Sectors    []BarSector
}

type BarSector struct {
	Color   string
	Portion float64
}

// Fight is a collection of reports about an encounter with a specific NPC
type Fight struct {
	// Id is a unique identifier for the fight
	Id        int

	// Target is the NPC that the fight is with
	Target    string

	// Reports collect all of the information involved in the fight
	Reports FightReportSet

	// StartTime is when we first noticed we were involved
	StartTime time.Time

	// LastActivity is the last time we noticed a log message involving this fight
	LastActivity time.Time
}
