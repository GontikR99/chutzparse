// +build wasm,web

package fight

import (
	"fmt"
	"github.com/gontikr99/chutzparse/internal/model/fight"
	"github.com/gontikr99/chutzparse/internal/model/iff"
	"github.com/gontikr99/chutzparse/internal/rpc"
	"github.com/gontikr99/chutzparse/internal/ui"
	"github.com/gontikr99/chutzparse/pkg/electron/ipc/ipcrenderer"
	"github.com/gontikr99/chutzparse/pkg/vuguutil"
	"github.com/vugu/vugu"
	"sort"
	"strconv"
	"time"
)

type Display struct {
	vuguutil.BackgroundComponent
	reportSet   fight.FightReportSet
	currentTab  string
	currentView vugu.Builder

	selectedLinkPet map[string]struct{}
	selectedLinkOwner map[string]struct{}
	selectedUnlinkPet map[string]struct{}
}

func (c *Display) Init(vCtx vugu.InitCtx) {
	c.currentTab = fight.ReportNames()[0]
	c.selectedUnlinkPet= map[string]struct{}{}
	c.selectedLinkPet= map[string]struct{}{}
	c.selectedLinkOwner= map[string]struct{}{}

	c.InitBackground(vCtx, c)
}

func (c *Display) RunInBackground() {
	newFight, newFightDone := listenForFights()
	defer newFightDone()
	petChange, petChangeDone := iff.ListenPets()
	defer petChangeDone()
	for {
		select {
		case <-c.Done():
			return
		case <-newFight:
			c.Env().Lock()
			c.Env().UnlockRender()
		case <-petChange:
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

type sboByLabel []ui.SelectBoxOption
func (s sboByLabel) Len() int {return len(s)}
func (s sboByLabel) Less(i, j int) bool {return s[i].Text < s[j].Text}
func (s sboByLabel) Swap(i, j int) {s[i],s[j] = s[j], s[i]}

func (c *Display) PetNames() []ui.SelectBoxOption {
	pets := iff.GetPets()
	for k, _ := range c.selectedUnlinkPet {
		if _, present := pets[k]; !present {
			delete(c.selectedUnlinkPet, k)
		}
	}
	if len(c.selectedUnlinkPet)==0 {
		c.selectedUnlinkPet[""]=struct{}{}
	}

	opts := []ui.SelectBoxOption{{"", ""}}
	for pet, owner := range iff.GetPets() {
		opts=append(opts, ui.SelectBoxOption{
			Text:  pet+" -> "+owner,
			Value: pet,
		})
	}
	sort.Sort(sboByLabel(opts))
	return opts
}

func (c *Display) UnlinkPet(event vugu.DOMEvent) {
	event.PreventDefault()
	event.StopPropagation()
	if len(c.selectedUnlinkPet)!=0 {
		for k, _ := range c.selectedUnlinkPet {
			if k!="" {
				go func() {
					rpc.UnlinkPet(ipcrenderer.Client, k)
				}()
				return
			}
		}
	}
}

func (c *Display) PotentialPetsOwners() []ui.SelectBoxOption {
	rawOpts:=map[string]struct{}{}
	for _, fight := range finishedFights {
		fight.Reports.Participants(rawOpts)
	}

	for k, _ := range c.selectedLinkPet {
		if _, present := rawOpts[k]; !present {
			delete(c.selectedLinkPet, k)
		}
	}
	if len(c.selectedLinkPet)==0 {
		c.selectedLinkPet[""]=struct{}{}
	}

	for k, _ := range c.selectedLinkOwner {
		if _, present := rawOpts[k]; !present {
			delete(c.selectedLinkOwner, k)
		}
	}
	if len(c.selectedLinkOwner)==0 {
		c.selectedLinkOwner[""]=struct{}{}
	}

	opts:=[]ui.SelectBoxOption{{"", ""}}
	pets := iff.GetPets()
	for rawOpt, _ := range rawOpts {
		if _, present := pets[rawOpt]; !present {
			opts = append(opts, ui.SelectBoxOption{
				Text:  rawOpt,
				Value: rawOpt,
			})
		}
	}
	sort.Sort(sboByLabel(opts))
	return opts
}

func (c *Display) LinkPet(event vugu.DOMEvent) {
	event.PreventDefault()
	event.StopPropagation()
	pet := ""
	for k, _ := range c.selectedLinkPet {
		pet=k
		break
	}
	owner := ""
	for k, _ := range c.selectedLinkOwner {
		owner=k
		break
	}
	if pet!="" && owner!="" {
		delete(c.selectedLinkPet, pet)
		delete(c.selectedLinkOwner, owner)
		c.selectedLinkPet[""]=struct{}{}
		c.selectedLinkOwner[""]=struct{}{}
		go func() {
			rpc.LinkPet(ipcrenderer.Client, pet, owner)
		}()
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
