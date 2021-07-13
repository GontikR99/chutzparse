package rpc

import "net/rpc"

type LookupSettingRequest struct {
	Key string
}

type LookupSettingResponse struct {
	Value   string
	Present bool
}

type SetSettingRequest struct {
	Key   string
	Value string
}

type SetSettingResponse struct{}

type ClearSettingRequest struct {
	Key string
}
type ClearSettingResponse struct{}

type SettingsServer interface {
	Lookup(key string) (value string, present bool, err error)
	Set(key string, value string) error
	Clear(key string) error
}

type StubSettings struct {
	settings SettingsServer
}

func (ss *StubSettings) LookupSetting(req *LookupSettingRequest, res *LookupSettingResponse) (err error) {
	res.Value, res.Present, err = ss.settings.Lookup(req.Key)
	return
}

func (ss *StubSettings) SetSetting(req *SetSettingRequest, res *SetSettingResponse) (err error) {
	err = ss.settings.Set(req.Key, req.Value)
	return
}

func (ss *StubSettings) ClearSetting(req *ClearSettingRequest, res *ClearSettingResponse) (err error) {
	err = ss.settings.Clear(req.Key)
	return
}

func LookupSetting(client *rpc.Client, key string) (string, bool, error) {
	req := &LookupSettingRequest{key}
	res := new(LookupSettingResponse)
	err := client.Call("StubSettings.LookupSetting", req, res)
	return res.Value, res.Present, err
}

func SetSetting(client *rpc.Client, key string, value string) error {
	req := &SetSettingRequest{Key: key, Value: value}
	res := new(SetSettingResponse)
	return client.Call("StubSettings.SetSetting", req, res)
}

func ClearSetting(client *rpc.Client, key string) error {
	return client.Call("StubSettings.ClearSetting", &ClearSettingRequest{key}, new(ClearSettingResponse))
}

func HandleSetting(setting SettingsServer) func(*rpc.Server) {
	ss := &StubSettings{setting}
	return func(server *rpc.Server) {
		server.Register(ss)
	}
}
