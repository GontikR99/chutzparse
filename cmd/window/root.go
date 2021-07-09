// +build wasm,web

package main

import (
	"github.com/gontikr99/chutzparse/cmd/window/fight"
	"github.com/gontikr99/chutzparse/cmd/window/welcome"
	"github.com/gontikr99/chutzparse/cmd/window/settings"
	"github.com/gontikr99/chutzparse/internal/place"
	"github.com/gontikr99/chutzparse/pkg/vuguutil"
	"github.com/vugu/vugu"
	"strings"
	"time"
)

type Root struct {
	vuguutil.BackgroundComponent
	LastPlace string
	Body      vugu.Builder
}

type routeEntry struct {
	Place       string
	DisplayName string
	ShowInNav   func() bool
	BodyGen     func() vugu.Builder
}

func (r routeEntry) ClassText() string {
	if place.GetPlace() == r.Place {
		return "nav-link active"
	} else {
		return "nav-link"
	}
}

var neverShow = func() bool { return false }
var alwaysShow = func() bool { return true }

var routes = []*routeEntry{
	{"", "Welcome", alwaysShow, func() vugu.Builder { return &welcome.Welcome{} }},
	{"fight", "Fights", alwaysShow, func() vugu.Builder { return &fight.Display{} }},
	{"settings", "Settings", alwaysShow, func() vugu.Builder { return &settings.Settings{} }},
}

func (c *Root) Init(vCtx vugu.InitCtx) {
	c.Body = &welcome.Welcome{}
	c.InitBackground(vCtx, c)
}

func (c *Root) RunInBackground() {
	lastPlace := place.GetPlace()
	for {
		<-time.After(10 * time.Millisecond)
		curPlace := place.GetPlace()
		if lastPlace != curPlace {
			c.Env().Lock()
			lastPlace = curPlace
			c.Env().UnlockRender()
		}
	}
}

func (c *Root) Compute(ctx vugu.ComputeCtx) {
	fullPlace := place.GetPlace()
	curPlace := strings.Split(fullPlace, ":")[0]
	if curPlace == c.LastPlace {
		return
	}
	for _, route := range routes {
		if route.Place == curPlace {
			c.Body = route.BodyGen()
			c.LastPlace = curPlace
			return
		}
	}
	place.NavigateTo(ctx.EventEnv(), "")
}
