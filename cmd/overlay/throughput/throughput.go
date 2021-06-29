package throughput

import (
	"fmt"
	"github.com/gontikr99/chutzparse/internal/parse_model/parsecomms"
	"github.com/gontikr99/chutzparse/internal/parse_model/parsedefs"
	"github.com/gontikr99/chutzparse/pkg/vuguutil"
	"github.com/vugu/vugu"
	"math"
)

// Throughput is the Vugu component representing the throughput display
type Throughput struct {
	vuguutil.BackgroundComponent
	AttrMap    vugu.AttrMap
	TopBars    []parsedefs.ThroughputBar
	BottomBars []parsedefs.ThroughputBar

	TextSize    string
	TextColor string
	TextElevation float64
	StepRadius  float64
	StepMargin float64
	InitialStep int
	ArcWidth    float64
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
	inChan := parsecomms.ThroughputListen(c.Ctx)
	for {
		select {
		case <-c.Done():
			return
		case barSet := <-inChan:
			c.Env().Lock()
			if len(barSet)==0 {
				c.TopBars = nil
				c.BottomBars = nil
			} else {
				c.TopBars = barSet[0].TopBars
				c.BottomBars = barSet[0].BottomBars
			}
			c.Env().UnlockRender()
		}
	}
}


func (c *Throughput) radius(step int) float64 {
	return (c.StepRadius+c.StepMargin)*float64(c.InitialStep+step)
}


// arcBase calculates coordinates for a portion of a bottom arc of one of the display bars.  Parameters are:
// * index: [-len(c.TopBars) .. -1] for the top bars, [0 .. len(c.TopBars)-1] for the bottom bars
// * arcStart, arcEnd: 0 for the leftmost point on the specified bar, 1 for the rightmost point.
func (c *Throughput) arcBase(barIdx int, arcStart float64, arcEnd float64) (radius, sx, sy, ex, ey float64, dirFlag int) {
	radius = (c.StepRadius + c.StepMargin) * float64(c.InitialStep)
	var startAngle, endAngle, yOffset float64
	if barIdx < 0 {
		startAngle = math.Pi/2 + c.ArcWidth/2 - arcStart*c.ArcWidth
		endAngle = math.Pi/2 + c.ArcWidth/2 - arcEnd*c.ArcWidth
		yOffset = -(c.StepRadius+c.StepMargin)*float64(-1-barIdx)
		dirFlag = 1
	} else {
		startAngle = 3*math.Pi/2 - c.ArcWidth/2 + arcStart*c.ArcWidth
		endAngle = 3*math.Pi/2 - c.ArcWidth/2 + arcEnd*c.ArcWidth
		yOffset = (c.StepRadius+c.StepMargin)*float64(barIdx+1)
		dirFlag = 0
	}
	sy = yOffset - radius*math.Sin(startAngle)
	ey = yOffset - radius*math.Sin(endAngle)

	sx = radius * math.Cos(startAngle)
	ex = radius * math.Cos(endAngle)
	return
}

// textPath calculates the text describing the arc that a bar's text will lay upon
func (c *Throughput) textPath(barIdx int) string {
	r, sx, sy, ex, ey, dir := c.arcBase(barIdx, 0.05, 0.95)
	return fmt.Sprintf("M %f %f A %f %f 0 0 %d %f %f", sx, sy-c.TextElevation, r, r, dir, ex, ey-c.TextElevation)
}

// sectorPath calculates the text describing how to draw a
func (c *Throughput) sectorPath(pathIdx int, arcStart, arcEnd float64) string {
	r, sx, sy, ex, ey, dir := c.arcBase(pathIdx, arcStart, arcEnd)
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
	index int
	bar   *parsedefs.ThroughputBar
}

// packBars combines the top and bottom bars into a single array for consumption by the web page
func (c *Throughput) packBars() []packedBar {
	var result []packedBar
	for pathIdx, _ := range c.TopBars {
		result = append(result, packedBar{-1-pathIdx, &c.TopBars[pathIdx]})
	}
	for pathIdx, _ := range c.BottomBars {
		result = append(result, packedBar{pathIdx, &c.BottomBars[pathIdx]})
	}
	return result
}


type packedBarSector struct {
	barIndex int
	color    string
	arcStart float64
	arcEnd   float64
}

// packSectors takes all the BarSlices in the top and bottom bars and packs them into a single array for consumption by
// the web page
func (c *Throughput) packSectors() []packedBarSector {
	var result []packedBarSector
	for _, bar := range c.packBars() {
		var accum float64
		for _, sector := range bar.bar.Sectors {
			if accum==0 {
				result = append(result, packedBarSector{
					barIndex: bar.index,
					color:    sector.Color,
					arcStart: 0.5-sector.Portion/2,
					arcEnd:   0.5+sector.Portion/2,
				})
			} else {
				result = append(result, packedBarSector{
					barIndex: bar.index,
					color:    sector.Color,
					arcStart: 0.5-accum/2-sector.Portion/2,
					arcEnd:   0.5-accum/2,
				})
				result = append(result, packedBarSector{
					barIndex: bar.index,
					color:    sector.Color,
					arcStart: 0.5+accum/2,
					arcEnd:   0.5+accum/2+sector.Portion/2,
				})
			}
			accum += sector.Portion
		}
	}
	return result
}