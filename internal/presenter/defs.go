package presenter

// HitEvent represents a single strike or heal that is shown on the hit display
type HitEvent struct {
	Text  string
	Color string
	Big   bool
}

const ChannelHitTop = "hitDisplayTop"
const ChannelHitBottom = "hitDisplayBottom"

type ThroughputState struct {
	FightId    int
	TopBars    []ThroughputBar
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

	LeftStyle   string
	CenterStyle string
	RightStyle  string

	Sectors []BarSector
}

type BarSector struct {
	Color   string
	Portion float64
}

const ThroughputBarCount = 10

const MonoStyle = "font-family: \"Anka Coder Narrow\";"
