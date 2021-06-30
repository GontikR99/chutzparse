package hit_display

import (
	"context"
	"fmt"
	"github.com/gontikr99/chutzparse/internal/parse_model/parsecomms"
	"github.com/gontikr99/chutzparse/internal/parse_model/parsedefs"
	"github.com/gontikr99/chutzparse/pkg/dom/document"
	"github.com/vugu/vugu"
	"math"
	"math/rand"
	"strings"
	"time"
)

// HitDisplay is the Vugu component which renders hit events
type HitDisplay struct {
	AttrMap   vugu.AttrMap
	allocated []int
}

func (c *HitDisplay) Run() {
	rand.Seed(time.Now().Unix())
	c.allocated = make([]int, pathCount)
	defsElem := document.GetElementById("hitDefs")
	for i:=0;i<pathCount;i++ {
		pathElem := document.CreateElementNS(namespaceSvg, "path")
		pathElem.SetAttribute("id", fmt.Sprintf("hit-path%d", i))
		pathElem.SetAttribute("fill", "none")
		pathElem.SetAttribute("d", textPath(i))
		defsElem.AppendChild(pathElem)
	}
	go func() {
		c.RunInBackground()
	}()
}

func (c *HitDisplay) RunInBackground() {
	topEvent := parsecomms.HitDisplayListen(context.Background(), parsedefs.ChannelHitTop)
	bottomEvent := parsecomms.HitDisplayListen(context.Background(), parsedefs.ChannelHitBottom)
	topSide:=0
	bottomSide:=0
	for {
		select {
		case hde := <- topEvent:
			c.drawRandomRange(hde.Text, hde.Color, hde.Big, (0+topSide)*pathCount/4, (1+topSide)*pathCount/4)
			topSide = 1 - topSide
		case hde:= <-bottomEvent:
			c.drawRandomRange(hde.Text, hde.Color, hde.Big, (2+bottomSide)*pathCount/4, (3+bottomSide)*pathCount/4)
			bottomSide = 1 - bottomSide
		}
	}
}


const pathCount = 32 // number of text paths to reserve
const screenSeconds = 3 // how long a hit text remains on screen

// drawRandomRange draws a given text string with a given color into one of the text paths of the hit display.
// The path to use is chosen randomly, within the range of [rangeBegin, rangeEnd).  If possible, an unused
// path is chosen.
//
// Text paths are numbered left to right in the range [0, pathCount/2) for the top, and [pathCount/2, pathCount) for the
// bottom.
func (c *HitDisplay) drawRandomRange(text string, color string, big bool, rangeBegin int, rangeEnd int) {
	options := []int{}
	for i := rangeBegin; i < rangeEnd; i++ {
		if c.allocated[(i+pathCount)%pathCount]==0 {
			options = append(options, i)
		}
	}
	if len(options) == 0 {
		c.draw(rangeBegin+rand.Intn(rangeEnd-rangeBegin), text, color, big)
	} else {
		c.draw(options[rand.Intn(len(options))], text, color, big)
	}
}

// tweakable drawing parameters
const outerBox = 160
const strokeSegments = 8
const centerRadius = 40
const skipYSteps = 5

// pathParams calculates the beginning and ending points for a specfied text path.
//
// Text paths are numbered left to right in the range [0, pathCount/2) for the top, and [pathCount/2, pathCount) for the
// bottom.
func pathParams(pathIndex int) (startX, startY, endX, endY float64) {
	var startAngle float64
	var yStep int
	if pathIndex < pathCount/4 {
		endX = -outerBox
		yStep = pathCount/4 - 1 - pathIndex
		startAngle = math.Pi - 2*math.Pi*float64(pathIndex+1)/(pathCount+4)
	} else if pathIndex < 2*pathCount/4 {
		endX = outerBox
		yStep = pathIndex - pathCount/4
		startAngle = math.Pi/2 - 2*math.Pi*float64(pathIndex-pathCount/4+1)/(pathCount+4)
	} else if pathIndex < 3*pathCount/4 {
		endX = -outerBox
		yStep = pathCount/4 + skipYSteps + pathIndex - 2*pathCount/4
		startAngle = math.Pi + 2*math.Pi*float64(pathIndex-2*pathCount/4+1)/(pathCount+4)
	} else {
		endX = outerBox
		yStep = pathCount/4 + skipYSteps + pathCount/4 - 1 - (pathIndex - 3*pathCount/4)
		startAngle = 2*math.Pi - 2*math.Pi*float64(pathCount-pathIndex)/(pathCount+4)
	}
	startX = centerRadius * math.Cos(startAngle)
	startY = -centerRadius * math.Sin(startAngle)
	endY = -outerBox + 2*outerBox*float64(yStep)/(pathCount/2+skipYSteps-1)
	return
}

// textPath returns the SVG path "d" parameter describing a specific text path.
//
// Text paths are numbered left to right in the range [0, pathCount/2) for the top, and [pathCount/2, pathCount) for the
// bottom.
func textPath(pathIndex int) string {
	startX, startY, endX, endY := pathParams(pathIndex)
	path := strings.Builder{}
	path.WriteString(fmt.Sprintf("M %f %f", startX, startY))
	for i := 0; i < strokeSegments+1; i++ {
		progress := float64(i) / strokeSegments
		pointX := (endX-startX)*progress + startX
		pointY := (endY-startY)*(1-(1-progress)*(1-progress)) + startY
		path.WriteString(fmt.Sprintf(" L %.1f %.1f", pointX, pointY))
	}
	return path.String()
}

var drawIdGen = 0 // provide each displayed text item a unique id
const namespaceSvg = "http://www.w3.org/2000/svg"

// draw places text upon a specified text path, starts it moving, and arranges it to
// be removed at the end of its animation.
func (c *HitDisplay) draw(pathIndex int, text string, color string, big bool) {
	c.allocated[pathIndex]++
	svgElem := document.GetElementById("hit")

	txtId := fmt.Sprintf("hit-text%d", drawIdGen)
	animId := fmt.Sprintf("hit-anim%d", drawIdGen)
	drawIdGen++

	txtElem := document.CreateElementNS(namespaceSvg, "text")
	txtElem.AppendChild(document.CreateTextNode(text))

	txtElem.SetAttribute("id", txtId)
	txtElem.SetAttribute("fill", color)
	txtElem.SetAttribute("fill-opacity", "90%")
	txtElem.SetAttribute("font-family", "Arial Black")
	if big {
		txtElem.SetAttribute("font-size", "20px")
	} else {
		txtElem.SetAttribute("font-size", "10px")
	}
	txtElem.SetAttribute("font-weight", "bolder")
	txtElem.SetAttribute("text-anchor", "middle")
	txtElem.SetAttribute("stroke", "black")
	txtElem.SetAttribute("stroke-width", "3")
	txtElem.SetAttribute("vector-effect", "non-scaling-stroke")
	startX, startY, _, _ := pathParams(pathIndex)
	txtElem.SetAttribute("x", fmt.Sprintf("%.1f", startX))
	txtElem.SetAttribute("y", fmt.Sprintf("%.1f", startY))
	svgElem.AppendChild(txtElem)

	animMotionElem := document.CreateElementNS(namespaceSvg, "animateMotion")
	animMotionElem.SetAttribute("id", animId)
	animMotionElem.SetAttribute("href", fmt.Sprintf("#%s", txtId))
	animMotionElem.SetAttribute("dur", fmt.Sprintf("%ds", screenSeconds))
	animMotionElem.SetAttribute("fill", "freeze")
	animMotionElem.SetAttribute("keyTimes", keyTimes)
	animMotionElem.SetAttribute("keyPoints", keyPoints)
	animMotionElem.SetAttribute("calcMode", "linear")
	animMotionElem.SetAttribute("begin", "click")
	animMotionElem.SetAttribute("repeatCount", "1")

	mpathElem := document.CreateElementNS(namespaceSvg, "mpath")
	mpathElem.SetAttribute("href", fmt.Sprintf("#hit-path%d", pathIndex))
	animMotionElem.AppendChild(mpathElem)
	svgElem.AppendChild(animMotionElem)

	go func() {
		// StartMain animation
		<-time.After(10 * time.Millisecond)
		txtElem.SetAttribute("x", 0)
		txtElem.SetAttribute("y", 0)
		animMotionElem.JSValue().Call("beginElement")

		// Remove/reclaim text
		<-time.After(screenSeconds * time.Second)
		animMotionElem.Remove()
		txtElem.Remove()
		c.allocated[pathIndex]--
	}()
}
var keyTimes = (func() string {
	sb := strings.Builder{}
	needSep := false
	for i := 0; i < strokeSegments+1; i++ {
		if needSep {
			sb.WriteString("; ")
		} else {
			needSep = true
		}
		progress := float64(i) / strokeSegments
		sb.WriteString(fmt.Sprintf("%.3f", progress))
	}
	return sb.String()
})()

var keyPoints = (func () string {
	sb := strings.Builder{}
	needSep := false
	for i := 0; i < strokeSegments+1; i++ {
		if needSep {
			sb.WriteString("; ")
		} else {
			needSep = true
		}
		progress := float64(i) / strokeSegments
		progress = 1 - math.Pow((1-progress), 4)
		sb.WriteString(fmt.Sprintf("%.3f", progress))
	}
	return sb.String()
})()