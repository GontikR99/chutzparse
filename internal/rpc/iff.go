package rpc

//go:generate ../../build/rpcgen iff.go

type IffControl interface {
	Unlink(pet string) error
	Link(pet string, owner string) error
}