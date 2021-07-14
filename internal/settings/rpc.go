package settings

//go:generate ../../build/rpcgen rpc.go

// Settings provides renderers the ability to query and change global program settings.
type Settings interface {
	Lookup(key string) (value string, present bool, err error)
	Set(key string, value string) error
	Clear(key string) error
}
