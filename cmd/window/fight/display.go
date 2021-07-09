// +build wasm,web

package fight

import (
	"fmt"
	"github.com/gontikr99/chutzparse/internal/model/fight"
	"github.com/gontikr99/chutzparse/internal/rpc"
	"github.com/gontikr99/chutzparse/internal/ui"
	"github.com/gontikr99/chutzparse/pkg/electron/ipc/ipcrenderer"
	"github.com/gontikr99/chutzparse/pkg/vuguutil"
	"github.com/vugu/vugu"
	"strconv"
	"time"
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
		duration := fgt.LastActivity.Sub(fgt.StartTime) / time.Second
		if duration < 0 {
			duration = 0
		}
		opts = append(opts, ui.SelectBoxOption{
			Text: fmt.Sprintf("[%02d:%02d:%02d +%4ds] %s",
				fgt.StartTime.Hour(), fgt.StartTime.Minute(), fgt.StartTime.Second(),
				duration,
				fgt.Target,
			),
			Value: strconv.FormatInt(int64(fgt.Id), 10),
		})
	}
	i := 0
	j := len(opts) - 1
	for i < j {
		opts[i], opts[j] = opts[j], opts[i]
		i++
		j--
	}
	return opts
}

func (c *Display) CopySummary(event vugu.DOMEvent) {
	event.PreventDefault()
	event.StopPropagation()
	if report, ok := c.reportSet[c.currentTab]; ok {
		text := report.Summarize()
		go func() {
			rpc.CopyClipboard(ipcrenderer.Client, text)
		}()
	}
}
