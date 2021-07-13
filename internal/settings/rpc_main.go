// +build wasm,electron

package settings

import "net/rpc"

type settingsServer struct{}

func (s settingsServer) Lookup(key string) (value string, present bool, err error) {
	return LookupSetting(key)
}

func (s settingsServer) Set(key string, value string) error {
	return SetSetting(key, value)
}

func (s settingsServer) Clear(key string) error {
	return ClearSetting(key)
}

func HandleRPC() func(server *rpc.Server) {
	return handleSettings(settingsServer{})
}