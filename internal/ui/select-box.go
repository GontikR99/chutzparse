// +build wasm,web

package ui

import (
	"fmt"
	"github.com/gontikr99/chutzparse/pkg/vuguutil"
	"github.com/vugu/vugu"
	"syscall/js"
)

var selectBoxIdGen int

type SelectBox struct {
	vuguutil.BackgroundComponent
	AttrMap  vugu.AttrMap
	Selected map[string]struct{}
	Change   SelectBoxChangeHandler
	Options  []SelectBoxOption

	idStr string
}

type SelectBoxOption struct {
	Text  string
	Value string
}

func (c *SelectBox) optNodeId(idx int) string {
	return fmt.Sprintf("%s-option-%d", c.idStr, idx)
}

func (c *SelectBox) Init(vCtx vugu.InitCtx) {
	selectBoxIdGen++
	c.idStr = fmt.Sprintf("selectbox-%d", selectBoxIdGen)
	if c.Selected == nil {
		c.Selected = make(map[string]struct{})
	}

	activeSelectBoxes[c.idStr] = c

	c.InitBackground(vCtx, c)
	c.ListenForRender()
}

func (c *SelectBox) RunInBackground() {
	defer func() {
		delete(activeSelectBoxes, c.idStr)
	}()
	for {
		select {
		case <-c.Done():
			return
		case <-c.Rendered():
			for idx, optVal := range c.Options {
				optElt := vuguutil.GetElementByNodeId(c.optNodeId(idx))
				if optElt == nil {
					continue
				}
				eltVal := optElt.GetAttribute("value")
				_, opSel := c.Selected[optVal.Value]
				if eltVal.IsNull() || eltVal.IsUndefined() || eltVal.String() != optVal.Value || !opSel {
					optElt.JSValue().Set("selected", false)
				} else {
					optElt.JSValue().Set("selected", true)
				}
			}
		}
	}
}

type selectBoxChangeEvent struct {
	selected map[string]struct{}
	env      vugu.EventEnv
	box      *SelectBox
}

var activeSelectBoxes = make(map[string]*SelectBox)

const sbChangeFunc = "cmiUiSelectBoxChange"

func init() {
	js.Global().Set(sbChangeFunc, js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		event := args[0]
		boxName := args[1].String()
		if box, ok := activeSelectBoxes[boxName]; ok {
			go func() {
				box.Env().Lock()
				box.onChange(vuguutil.NewVuguEvent(event, box.Env()))
				box.Env().UnlockRender()
			}()
		}
		return nil
	}))
}

func (c *SelectBox) onChangeHookText() string {
	return sbChangeFunc + "(event, \"" + c.idStr + "\")"
}

func (s *selectBoxChangeEvent) Selected() map[string]struct{} { return s.selected }
func (s *selectBoxChangeEvent) Env() vugu.EventEnv            { return s.box.Env() }

//func (s *selectBoxChangeEvent) SetSelected(s2 string) { s.box.Value = s2 }

func (c *SelectBox) onChange(event vugu.DOMEvent) {
	oldValue := make(map[string]struct{})
	for k, v := range c.Selected {
		oldValue[k] = v
	}

	changeDetected := false
	for idx, optVal := range c.Options {
		optElt := vuguutil.GetElementByNodeId(c.optNodeId(idx))
		if optElt == nil {
			continue
		}
		eltVal := optElt.GetAttribute("value")
		isSelected := optElt.JSValue().Get("selected").Bool()
		_, wasSelected := oldValue[optVal.Value]
		if eltVal.IsNull() || eltVal.IsUndefined() || eltVal.String() != optVal.Value || !isSelected {
			changeDetected = changeDetected || wasSelected
			delete(c.Selected, optVal.Value)
		} else {
			changeDetected = changeDetected || !wasSelected
			c.Selected[optVal.Value] = struct{}{}
		}
	}
	selectionCopy := make(map[string]struct{})
	for k, v := range c.Selected {
		selectionCopy[k] = v
	}
	if c.Change != nil && changeDetected {
		c.Change.SelectBoxChangeHandle(&selectBoxChangeEvent{
			selected: selectionCopy,
			box:      c,
			env:      c.Env(),
		})
	}
}
