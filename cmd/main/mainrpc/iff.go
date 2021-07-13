// +build wasm,electron

package mainrpc

import (
	"github.com/gontikr99/chutzparse/internal/iff"
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
	register(iff.HandleControl(iffStub{}))
}