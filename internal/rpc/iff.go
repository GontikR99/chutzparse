package rpc

//go:generate ../../build/rpcgen iff.go

// IffControl instructs the main process to update some IFF status.
type IffControl interface {
	Unlink(pet string) error
	Link(pet string, owner string) error
}