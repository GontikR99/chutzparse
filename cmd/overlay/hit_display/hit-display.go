package hit_display

import (
	"fmt"
	"github.com/gontikr99/chutzparse/internal/dmodel"
	"github.com/gontikr99/chutzparse/pkg/dom/document"
	"github.com/vugu/vugu"
	"math"
	"math/rand"
	"strings"
	"time"
)

const strokeSegments = 8
const centerRadius = 80
const pathCount = 32
const skipYSteps = 5
const screenSeconds = 3
const svgNS = "http://www.w3.org/2000/svg"

type HitDisplay struct {
	AttrMap   vugu.AttrMap
	allocated []int
}

func (c *HitDisplay) Init(ctx vugu.InitCtx) {
	rand.Seed(time.Now().Unix())
	c.allocated = make([]int, pathCount)
	dmodel.HitDisplayListen(dmodel.ChannelTopTarget, func(evt *dmodel.HitDisplayEvent) {
		c.drawRandomRange(evt.Text, evt.Color, 0, pathCount/2)
	})
	dmodel.HitDisplayListen(dmodel.ChannelBottomTarget, func(evt *dmodel.HitDisplayEvent) {
		c.drawRandomRange(evt.Text, evt.Color, pathCount/2, pathCount)
	})
}

func (c *HitDisplay) drawRandomRange(text string, color string, rangeBegin int, rangeEnd int) {
	options := []int{}
	for i := rangeBegin; i < rangeEnd; i++ {
		if c.allocated[(i+pathCount)%pathCount]==0 {
			options = append(options, i)
		}
	}
	if len(options) == 0 {
		c.draw(rangeBegin+rand.Intn(rangeEnd-rangeBegin), text, color)
	} else {
		c.draw(options[rand.Intn(len(options))], text, color)
	}
}

func pathParams(pathIndex int) (startX, startY, endX, endY float64) {
	var startAngle float64
	var yStep int
	if pathIndex < pathCount/4 {
		endX = -180
		yStep = pathCount/4 - 1 - pathIndex
		startAngle = math.Pi - 2*math.Pi*float64(pathIndex+1)/(pathCount+4)
	} else if pathIndex < 2*pathCount/4 {
		endX = 180
		yStep = pathIndex - pathCount/4
		startAngle = math.Pi/2 - 2*math.Pi*float64(pathIndex-pathCount/4+1)/(pathCount+4)
	} else if pathIndex < 3*pathCount/4 {
		endX = -180
		yStep = pathCount/4 + skipYSteps + pathIndex - 2*pathCount/4
		startAngle = math.Pi + 2*math.Pi*float64(pathIndex-2*pathCount/4+1)/(pathCount+4)
	} else {
		endX = 180
		yStep = pathCount/4 + skipYSteps + pathCount/4 - 1 - (pathIndex - 3*pathCount/4)
		startAngle = 2*math.Pi - 2*math.Pi*float64(pathCount-pathIndex)/(pathCount+4)
	}
	startX = centerRadius * math.Cos(startAngle)
	startY = -centerRadius * math.Sin(startAngle)
	endY = -180 + 360*float64(yStep)/(pathCount/2+skipYSteps-1)
	return
}

func textPath(barIndex int) string {
	startX, startY, endX, endY := pathParams(barIndex)
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

var drawIdGen = 0

func (c *HitDisplay) draw(pathIndex int, text string, color string) {
	pathIndex = ((pathIndex % pathCount) + pathCount) % pathCount
	c.allocated[pathIndex]++
	svgElem := document.GetElementById("hit")

	txtId := fmt.Sprintf("hit-text%d", drawIdGen)
	animId := fmt.Sprintf("hit-anim%d", drawIdGen)
	drawIdGen++

	txtElem := document.CreateElementNS(svgNS, "text")
	txtElem.AppendChild(document.CreateTextNode(text))

	txtElem.SetAttribute("id", txtId)
	txtElem.SetAttribute("fill", color)
	txtElem.SetAttribute("fill-opacity", "90%")
	txtElem.SetAttribute("font-size", "12px")
	txtElem.SetAttribute("font-weight", "bolder")
	txtElem.SetAttribute("text-anchor", "middle")
	txtElem.SetAttribute("stroke", "none")
	startX, startY, _, _ := pathParams(pathIndex)
	txtElem.SetAttribute("x", fmt.Sprintf("%.1f", startX))
	txtElem.SetAttribute("y", fmt.Sprintf("%.1f", startY))
	svgElem.AppendChild(txtElem)

	animMotionElem := document.CreateElementNS(svgNS, "animateMotion")
	animMotionElem.SetAttribute("id", animId)
	animMotionElem.SetAttribute("href", fmt.Sprintf("#%s", txtId))
	animMotionElem.SetAttribute("dur", fmt.Sprintf("%ds", screenSeconds))
	animMotionElem.SetAttribute("fill", "freeze")
	animMotionElem.SetAttribute("keyTimes", keyTimes)
	animMotionElem.SetAttribute("keyPoints", keyPoints)
	animMotionElem.SetAttribute("calcMode", "linear")
	animMotionElem.SetAttribute("begin", "click")
	animMotionElem.SetAttribute("repeatCount", "1")

	mpathElem := document.CreateElementNS(svgNS, "mpath")
	mpathElem.SetAttribute("href", fmt.Sprintf("#hit-path%d", pathIndex))
	animMotionElem.AppendChild(mpathElem)
	svgElem.AppendChild(animMotionElem)

	go func() {
		<-time.After(10 * time.Millisecond)
		txtElem.SetAttribute("x", 0)
		txtElem.SetAttribute("y", 0)
		document.GetElementById(animId).JSValue().Call("beginElement")
	}()

	go func() {
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