package randoms

//go:generate ../../build/rpcgen randomsrpc.go

type Randoms interface {
	Reset() error
	Cull() error
	FetchRandoms() ([]*RollGroup, error)
}
