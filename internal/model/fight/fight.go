package fight

import (
	"time"
)

// Fight is a collection of reports about an encounter with a specific NPC
type Fight struct {
	// Id is a unique identifier for the fight
	Id int

	// Target is the NPC that the fight is with
	Target string

	// Reports collect all of the information involved in the fight
	Reports FightReportSet

	// StartTime is when we first noticed we were involved
	StartTime time.Time

	// LastActivity is the last time we noticed a log message involving this fight
	LastActivity time.Time
}

const ChannelFinishedFights = "finishedFights"
