// +build wasm,web

package settings

import (
	"github.com/gontikr99/chutzparse/internal/eqspec"
	"github.com/gontikr99/chutzparse/internal/settings"
	"github.com/gontikr99/chutzparse/internal/ui"
	"github.com/gontikr99/chutzparse/pkg/console"
	"github.com/gontikr99/chutzparse/pkg/electron/ipc/ipcrenderer"
	"github.com/gontikr99/chutzparse/pkg/vuguutil"
	"github.com/vugu/vugu"
)

type Settings struct {
	vuguutil.BackgroundComponent
	EqDir            *ConfiguredValue
	LinkObviousPets  bool
	EnableFlyingText bool
	EnableMeters     bool
}

var scanctl = eqspec.NewScanControlClient(ipcrenderer.Client)

func (c *Settings) Init(vCtx vugu.InitCtx) {
	c.EqDir = &ConfiguredValue{
		Key:      settings.EverQuestDirectory,
		Callback: func(s string) { scanctl.Restart() },
	}

	c.EqDir.Init(vCtx)
	c.InitBackground(vCtx, c)
}

func (c *Settings) RunInBackground() {
	sc, doneFunc := ipcrenderer.Endpoint{}.Listen(settings.ChannelChange)
	defer doneFunc()
	for {
		go c.refreshCheckbox(c.Env(), settings.LinkObviousPets, &c.LinkObviousPets)
		go c.refreshCheckbox(c.Env(), settings.ShowFlyingHits, &c.EnableFlyingText)
		go c.refreshCheckbox(c.Env(), settings.ShowMeters, &c.EnableMeters)
		select {
		case <-c.Done():
			return
		case <-sc:
		}
	}
}

var rpcset = settings.NewSettingsClient(ipcrenderer.Client)

func (c *Settings) refreshCheckbox(env vugu.EventEnv, cbName string, cbValue *bool) {
	value, present, err := rpcset.Lookup(cbName)
	if present && err == nil {
		env.Lock()
		if value == "true" {
			*cbValue = true
		} else {
			*cbValue = false
		}
		env.UnlockRender()
	}
}

func (c *Settings) ToggleCheckbox(event vugu.DOMEvent, settingName string, settingValue *bool) {
	*settingValue = event.JSEventTarget().Get("checked").Truthy()
	go func() {
		if *settingValue {
			rpcset.Set(settingName, "true")
		} else {
			rpcset.Set(settingName, "false")
		}
		c.refreshCheckbox(event.EventEnv(), settingName, settingValue)
	}()
}

var dirdlg = ui.NewDirectoryDialogClient(ipcrenderer.Client)

func (c *Settings) BrowseEqDir(event vugu.DOMEvent) {
	event.PreventDefault()
	go func() {
		newDir, err := dirdlg.Choose(c.EqDir.Value)
		if err != nil {
			console.Log(err.Error())
			return
		}
		event.EventEnv().Lock()
		c.EqDir.SetStringValue(newDir)
		event.EventEnv().UnlockRender()
	}()
}

type ConfiguredValue struct {
	Key      string
	Value    string
	Callback func(value string)
}

func (cv *ConfiguredValue) Init(ctx vugu.InitCtx) {
	go func() {
		value, present, err := rpcset.Lookup(cv.Key)
		if err == nil && present {
			ctx.EventEnv().Lock()
			cv.Value = value
			ctx.EventEnv().UnlockRender()
		}
	}()
}

func (cv *ConfiguredValue) StringValue() string {
	return cv.Value
}

func (cv *ConfiguredValue) SetStringValue(s string) {
	cv.Value = s
	go func() {
		rpcset.Set(cv.Key, s)
		if cv.Callback != nil {
			cv.Callback(s)
		}
	}()
}
