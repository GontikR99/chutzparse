// +build wasm,electron

package mainrpc

import (
	"github.com/gontikr99/chutzparse/internal/model/iff"
	"github.com/gontikr99/chutzparse/internal/rpc"
)

type iffStub struct{}

func (i iffStub) Unlink(pet string) error {
	iff.UnlinkPet(pet)
	return nil
}

func (i iffStub) Link(pet string, owner string) error {
	iff.MakePet(pet, owner)
	return nil
}

func init() {
	register(rpc.HandleIff(iffStub{}))
}