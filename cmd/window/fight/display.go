// +build wasm,web

package fight

import (
	"github.com/gontikr99/chutzparse/internal/model/fight"
	"github.com/gontikr99/chutzparse/internal/ui"
	"github.com/gontikr99/chutzparse/pkg/vuguutil"
	"github.com/vugu/vugu"
	"strconv"
)

type Display struct {
	vuguutil.BackgroundComponent
	reportSet   fight.FightReportSet
	currentTab  string
	currentView vugu.Builder
}

func (c *Display) Init(vCtx vugu.InitCtx) {
	c.currentTab = fight.ReportNames()[0]
	c.InitBackground(vCtx, c)
}

func (c *Display) RunInBackground() {
	newFight, newFightDone := listenForFights()
	defer newFightDone()
	for {
		select {
		case <-c.Done():
			return
		case <-newFight:
			c.Env().Lock()
			c.Env().UnlockRender()
		}
	}
}

func (c *Display) SetTab(event vugu.DOMEvent, tabName string) {
	event.PreventDefault()
	if c.currentTab != tabName {
		c.currentTab = tabName
		c.rebuildView()
	}
}

func (c *Display) SelectBoxChangeHandle(event ui.SelectBoxChangeEvent) {
	var activeReports []fight.FightReportSet
	for _, fight := range finishedFights {
		if _, present := event.Selected()[strconv.FormatInt(int64(fight.Id), 10)]; present {
			activeReports = append(activeReports, fight.Reports)
		}
	}
	if len(activeReports) == 0 {
		c.reportSet = nil
	} else if len(activeReports) == 1 {
		c.reportSet = activeReports[0]
	} else {
		c.reportSet = fight.MergeFightReports(activeReports)
	}
	c.rebuildView()
}

func (c *Display) rebuildView() {
	c.currentView = nil
	if c.reportSet == nil {
		return
	}
	if report, present := c.reportSet[c.currentTab]; present {
		c.currentView = report.Detail()
	}
}

func (c *Display) TabClass(tabName string) string {
	if c.currentTab == tabName {
		return "nav-link active"
	} else {
		return "nav-link"
	}
}

func (c *Display) FightNames() []ui.SelectBoxOption {
	var opts []ui.SelectBoxOption
	for _, fgt := range finishedFights {
		opts = append(opts, ui.SelectBoxOption{
			Text:  fgt.Target,
			Value: strconv.FormatInt(int64(fgt.Id), 10),
		})
	}
	return opts
}
