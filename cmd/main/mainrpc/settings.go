// +build wasm,electron

package mainrpc

import (
	"github.com/gontikr99/chutzparse/internal/settings"
)

type settingsServer struct{}

func (s settingsServer) Lookup(key string) (value string, present bool, err error) {
	return settings.LookupSetting(key)
}

func (s settingsServer) Set(key string, value string) error {
	return settings.SetSetting(key, value)
}

func (s settingsServer) Clear(key string) error {
	return settings.ClearSetting(key)
}

func init() {
	register(settings.HandleSettings(settingsServer{}))
}
