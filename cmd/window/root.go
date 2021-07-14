// +build wasm,web

package main

import (
	"github.com/gontikr99/chutzparse/cmd/window/fight"
	"github.com/gontikr99/chutzparse/cmd/window/settings"
	"github.com/gontikr99/chutzparse/cmd/window/welcome"
	"github.com/gontikr99/chutzparse/pkg/vuguutil"
	"github.com/vugu/vugu"
	"log"
	"net/url"
	"strings"
	"syscall/js"
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
	if GetPlace() == r.Place {
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
	lastPlace := GetPlace()
	for {
		<-time.After(10 * time.Millisecond)
		curPlace := GetPlace()
		if lastPlace != curPlace {
			c.Env().Lock()
			lastPlace = curPlace
			c.Env().UnlockRender()
		}
	}
}

func (c *Root) Compute(ctx vugu.ComputeCtx) {
	fullPlace := GetPlace()
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
	NavigateTo(ctx.EventEnv(), "")
}

var window = js.Global().Get("window")

func GetPlace() string {
	href := window.Get("location").Get("href").String()
	parsed, err := url.Parse(href)
	if err != nil {
		log.Println(err)
		return ""
	}
	return parsed.Fragment
}

func NavigateTo(env vugu.EventEnv, place string) {
	go func() {
		env.Lock()
		window.Get("history").Call("pushState", nil, "", "#"+place)
		env.UnlockRender()
	}()
}
