package throughput

import (
	"fmt"
	"github.com/gontikr99/chutzparse/internal/presenter"
	"github.com/gontikr99/chutzparse/pkg/vuguutil"
	"github.com/vugu/vugu"
	"math"
)

// Throughput is the Vugu component representing the throughput display
type Throughput struct {
	vuguutil.BackgroundComponent
	AttrMap  vugu.AttrMap
	Displays []presenter.ThroughputState

	TextSize      string
	TextColor     string
	TextElevation float64
	StepRadius    float64
	StepMargin    float64
	InitialStep   int
	ArcWidth      float64
}

const defaultTextSize = "8px"
const defaultTextColor = "yellow"
const defaultTextElevation = 2
const defaultStepRadius = 12
const defaultStepMargin = 2
const defaultInitialStep = 8
const defaultArcWidth = 2 * math.Pi / 3

func (c *Throughput) Init(vCtx vugu.InitCtx) {
	if c.TextSize == "" {
		c.TextSize = defaultTextSize
	}
	if c.TextColor == "" {
		c.TextColor = defaultTextColor
	}
	if c.TextElevation == 0 {
		c.TextElevation = defaultTextElevation
	}
	if c.StepRadius == 0 {
		c.StepRadius = defaultStepRadius
	}
	if c.StepMargin == 0 {
		c.StepMargin = defaultStepMargin
	}
	if c.InitialStep == 0 {
		c.InitialStep = defaultInitialStep
	}
	if c.ArcWidth == 0 {
		c.ArcWidth = defaultArcWidth
	}
	c.InitBackground(vCtx, c)
}

func (c *Throughput) RunInBackground() {
	inChan := presenter.ThroughputListen(c.Ctx)
	for {
		select {
		case <-c.Done():
			return
		case barSet := <-inChan:
			c.Env().Lock()
			c.Displays = barSet
			if len(c.Displays) > 3 {
				c.Displays = c.Displays[:3]
			}
			c.Env().UnlockRender()
		}
	}
}

func (c *Throughput) radius(step int) float64 {
	return (c.StepRadius + c.StepMargin) * float64(c.InitialStep+step)
}

// arcBase calculates coordinates for a portion of a bottom arc of one of the display bars.  Parameters are:
// * index: [-len(c.TopBars) .. -1] for the top bars, [0 .. len(c.TopBars)-1] for the bottom bars
// * arcStart, arcEnd: 0 for the leftmost point on the specified bar, 1 for the rightmost point.
func (c *Throughput) arcBase(displayIdx int, barIdx int, arcStart float64, arcEnd float64) (radius, sx, sy, ex, ey float64, dirFlag int) {
	radius = (c.StepRadius + c.StepMargin) * float64(c.InitialStep)
	xOffset := 1.1 * float64(displayIdx) * radius * math.Cos(math.Pi/2-c.ArcWidth/2)
	var startAngle, endAngle, yOffset float64
	if barIdx < 0 {
		startAngle = math.Pi/2 + c.ArcWidth/2 - arcStart*c.ArcWidth
		endAngle = math.Pi/2 + c.ArcWidth/2 - arcEnd*c.ArcWidth
		yOffset = -(c.StepRadius + c.StepMargin) * float64(-1-barIdx)
		dirFlag = 1
	} else {
		startAngle = 3*math.Pi/2 - c.ArcWidth/2 + arcStart*c.ArcWidth
		endAngle = 3*math.Pi/2 - c.ArcWidth/2 + arcEnd*c.ArcWidth
		yOffset = (c.StepRadius + c.StepMargin) * float64(barIdx+1)
		dirFlag = 0
	}
	sy = yOffset - radius*math.Sin(startAngle)
	ey = yOffset - radius*math.Sin(endAngle)

	sx = radius*math.Cos(startAngle) + xOffset
	ex = radius*math.Cos(endAngle) + xOffset
	return
}

// textPath calculates the text describing the arc that a bar's text will lay upon
func (c *Throughput) textPath(displayIdx int, barIdx int) string {
	r, sx, sy, ex, ey, dir := c.arcBase(displayIdx, barIdx, 0.05, 0.95)
	return fmt.Sprintf("M %f %f A %f %f 0 0 %d %f %f", sx, sy-c.TextElevation, r, r, dir, ex, ey-c.TextElevation)
}

// sectorPath calculates the text describing how to draw a
func (c *Throughput) sectorPath(displayIdx int, pathIdx int, arcStart, arcEnd float64) string {
	r, sx, sy, ex, ey, dir := c.arcBase(displayIdx, pathIdx, arcStart, arcEnd)
	return fmt.Sprintf(
		"M %f %f A %f %f 0 0 %d %f %f "+
			"l 0 %f "+
			"A %f %f 0 0 %d %f %f "+
			"Z",
		sx, sy, r, r, dir, ex, ey,
		-c.StepRadius,
		r, r, 1-dir, sx, sy-c.StepRadius)
}

type packedBar struct {
	display int
	index   int
	bar     *presenter.ThroughputBar
}

func (c *Throughput) displaysOffset() int {
	switch len(c.Displays) {
	case 1:
		return 0
	case 2:
		return -1
	case 3:
		return -2
	default:
		return 0
	}
}

// packBars combines the top and bottom bars into a single array for consumption by the web page
func (c *Throughput) packBars() []packedBar {
	var result []packedBar
	for displayNum, display := range c.Displays {
		for pathIdx, _ := range display.TopBars {
			result = append(result, packedBar{2*displayNum + c.displaysOffset(), -1 - pathIdx, &c.Displays[displayNum].TopBars[pathIdx]})
		}
		for pathIdx, _ := range display.BottomBars {
			result = append(result, packedBar{2*displayNum + c.displaysOffset(), pathIdx, &c.Displays[displayNum].BottomBars[pathIdx]})
		}
	}
	return result
}

type packedBarSector struct {
	displayIndex int
	barIndex     int
	color        string
	arcStart     float64
	arcEnd       float64
}

// packSectors takes all the BarSlices in the top and bottom bars and packs them into a single array for consumption by
// the web page
func (c *Throughput) packSectors() []packedBarSector {
	var result []packedBarSector
	for _, bar := range c.packBars() {
		var total float64
		for _, sector := range bar.bar.Sectors {
			total += sector.Portion
		}
		var accum float64
		for _, sector := range bar.bar.Sectors {
			result = append(result, packedBarSector{
				displayIndex: bar.display,
				barIndex:     bar.index,
				color:        sector.Color,
				arcStart:     (1.0-total)/2 + accum,
				arcEnd:       (1.0-total)/2 + accum + sector.Portion,
			})
			accum += sector.Portion
		}
	}
	return result
}
