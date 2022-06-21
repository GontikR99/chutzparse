//go:build wasm && web
// +build wasm,web

package randoms

import (
	"fmt"
	"github.com/gontikr99/chutzparse/internal/randoms"
	"github.com/gontikr99/chutzparse/pkg/electron/ipc/ipcrenderer"
	"github.com/gontikr99/chutzparse/pkg/vuguutil"
	"github.com/vugu/vugu"
)

var randomsctl = randoms.NewRandomsClient(ipcrenderer.Client)

type Randoms struct {
	vuguutil.BackgroundComponent
	CurrentRandoms []*randoms.RollGroup
}

func (c *Randoms) Init(vCtx vugu.InitCtx) {
	go func() {
		cr, _ := randomsctl.FetchRandoms()
		vCtx.EventEnv().Lock()
		c.CurrentRandoms = cr
		vCtx.EventEnv().UnlockRender()
	}()
	c.InitBackground(vCtx, c)
}

func (c *Randoms) RunInBackground() {
	sc, doneFunc := ipcrenderer.Endpoint{}.Listen(randoms.ChannelChange)
	defer doneFunc()
	for {
		select {
		case <-c.Done():
			return
		case <-sc:
			cr, _ := randomsctl.FetchRandoms()
			c.Env().Lock()
			c.CurrentRandoms = cr
			c.Env().UnlockRender()
		}
	}
}

func (c *Randoms) Reset(event vugu.DOMEvent) {
	go func() {
		randomsctl.Reset()
	}()
}

func title(rg *randoms.RollGroup) string {
	return fmt.Sprintf("%v .. %v", rg.Min, rg.Max)
}

func rvalue(roll *randoms.CharacterRoll) string {
	return fmt.Sprint(roll.Value)
}
