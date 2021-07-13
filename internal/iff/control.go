package iff

//go:generate ../../build/rpcgen control.go

// IffControl instructs the main process to update some IFF status.
type Control interface {
	Unlink(pet string) error
	Link(pet string, owner string) error
}
