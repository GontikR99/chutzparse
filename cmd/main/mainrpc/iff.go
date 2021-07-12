// +build wasm,electron

package mainrpc

import (
	iff2 "github.com/gontikr99/chutzparse/internal/iff"
	"github.com/gontikr99/chutzparse/internal/rpc"
)

type iffStub struct{}

func (i iffStub) Unlink(pet string) error {
	iff2.UnlinkPet(pet)
	return nil
}

func (i iffStub) Link(pet string, owner string) error {
	iff2.MakePet(pet, owner)
	return nil
}

func init() {
	register(rpc.HandleIff(iffStub{}))
}