//go:build wasm && web
// +build wasm,web

package bids

import (
	"github.com/gontikr99/chutzparse/internal/bids"
	"github.com/gontikr99/chutzparse/internal/settings"
	"github.com/gontikr99/chutzparse/pkg/electron/ipc/ipcrenderer"
	"github.com/gontikr99/chutzparse/pkg/vuguutil"
	"github.com/vugu/vugu"
	"strconv"
)

var rpcset = settings.NewSettingsClient(ipcrenderer.Client)
var bidsctl = bids.NewBidsClient(ipcrenderer.Client)

type Bids struct {
	vuguutil.BackgroundComponent
	CurrentBids []*bids.ItemBids
	StartCmd    string
	EndCmd      string
	Active      bool
	HasDump     bool
}

func (c *Bids) Init(vCtx vugu.InitCtx) {
	go func() {
		c.StartCmd, _, _ = rpcset.Lookup(settings.BidStartCmd)
		c.EndCmd, _, _ = rpcset.Lookup(settings.BidEndCmd)
		c.CurrentBids, _ = bidsctl.FetchBids()
		c.Active, _ = bidsctl.AuctionActive()
		c.HasDump, _ = bidsctl.HasGuildDump()

		vCtx.EventEnv().Lock()
		vCtx.EventEnv().UnlockRender()
	}()
	c.InitBackground(vCtx, c)
}

func (c *Bids) RunInBackground() {
	sc, doneFunc := ipcrenderer.Endpoint{}.Listen(bids.ChannelChange)
	defer doneFunc()
	for {
		select {
		case <-c.Done():
			return
		case <-sc:
			go func() {
				cb, _ := bidsctl.FetchBids()
				ac, _ := bidsctl.AuctionActive()
				hd, _ := bidsctl.HasGuildDump()
				c.Env().Lock()
				c.CurrentBids = cb
				c.Active = ac
				c.HasDump = hd
				c.Env().UnlockRender()
			}()
		}
	}
}

func bidDKP(bid int32) string {
	if bid < 0 {
		return "???"
	} else {
		return strconv.Itoa(int(bid))
	}
}

func (c *Bids) StartAuction() {
	go func() {
		bidsctl.Start()
	}()
}

func (c *Bids) EndAuction() {
	go func() {
		bidsctl.End()
	}()
}
