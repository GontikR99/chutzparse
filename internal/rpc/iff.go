package rpc

import "net/rpc"

type IffServer interface {
	Unlink(pet string) error
	Link(pet string, owner string) error
}

type StubIff struct {
	is IffServer
}

func (ic *StubIff) Unlink(req *UnlinkPetRequest, res *UnlinkPetResponse) error {
	return ic.is.Unlink(req.Pet)
}

func UnlinkPet(client *rpc.Client, pet string) error {
	req := &UnlinkPetRequest{pet}
	res := new(UnlinkPetResponse)
	return client.Call("StubIff.Unlink", req, res)
}

func (ic *StubIff) Link(req *LinkPetRequest, res *LinkPetResponse) error {
	return ic.is.Link(req.Pet, req.Owner)
}

func LinkPet(client *rpc.Client, pet string, owner string) error {
	req := &LinkPetRequest{pet, owner}
	res := new(LinkPetResponse)
	return client.Call("StubIff.Link", req, res)
}

type UnlinkPetRequest struct {
	Pet string
}

type UnlinkPetResponse struct {}

type LinkPetRequest struct {
	Pet string
	Owner string
}
type LinkPetResponse struct {}

func HandleIff(iff IffServer) func(server *rpc.Server) {
	ic := &StubIff{iff}
	return func(server *rpc.Server) {
		server.Register(ic)
	}
}