// +build wasm,web

package welcome

import (
	"github.com/gontikr99/chutzparse/internal/rpc"
	"github.com/gontikr99/chutzparse/internal/settings"
	"github.com/gontikr99/chutzparse/pkg/console"
	"github.com/gontikr99/chutzparse/pkg/electron/ipc/ipcrenderer"
	"github.com/vugu/vugu"
)

type Welcome struct {
	EqDir *ConfiguredValue
}

func (c *Welcome) Init(vCtx vugu.InitCtx) {
	c.EqDir = &ConfiguredValue{
		Key:      settings.EverQuestDirectory,
		Callback: func(s string) { rpc.RestartScan(ipcrenderer.Client) },
	}

	c.EqDir.Init(vCtx)
}

func (c *Welcome) BrowseEqDir(event vugu.DOMEvent) {
	event.PreventDefault()
	go func() {
		newDir, err := rpc.DirectoryDialog(ipcrenderer.Client, c.EqDir.Value)
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
		value, present, err := rpc.LookupSetting(ipcrenderer.Client, cv.Key)
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
		rpc.SetSetting(ipcrenderer.Client, cv.Key, s)
		if cv.Callback != nil {
			cv.Callback(s)
		}
	}()
}
