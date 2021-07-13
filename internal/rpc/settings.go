package rpc

//go:generate ../../build/rpcgen settings.go

// Settings provides renderers the ability to query and change global program settings.
type Settings interface {
	Lookup(key string) (value string, present bool, err error)
	Set(key string, value string) error
	Clear(key string) error
}
