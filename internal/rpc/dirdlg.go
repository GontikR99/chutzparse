package rpc

import "net/rpc"

type DirectoryDialogRequest struct {
	InitialDirectory string
}

type DirectoryDialogResponse struct {
	ChosenDirectory string
}

type StubDirectoryDialog struct {
	callback func(initial string) (chosen string, err error)
}

func (d *StubDirectoryDialog) DirectoryDialog(req *DirectoryDialogRequest, res *DirectoryDialogResponse) error {
	chosen, err := d.callback(req.InitialDirectory)
	*res = DirectoryDialogResponse{chosen}
	return err
}

func DirectoryDialog(client *rpc.Client, initial string) (string, error) {
	req := &DirectoryDialogRequest{initial}
	res := &DirectoryDialogResponse{}
	err := client.Call("StubDirectoryDialog.DirectoryDialog", req, res)
	return res.ChosenDirectory, err
}

func HandleDirectoryDialog(callback func(string) (string, error)) func(*rpc.Server) {
	dds := &StubDirectoryDialog{callback}
	return func(server *rpc.Server) {
		server.Register(dds)
	}
}
