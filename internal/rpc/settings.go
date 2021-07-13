package rpc

//go:generate ../../build/rpcgen settings.go

type Settings interface {
	Lookup(key string) (value string, present bool, err error)
	Set(key string, value string) error
	Clear(key string) error
}
