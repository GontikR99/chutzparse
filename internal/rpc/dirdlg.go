package rpc

//go:generate ../../build/rpcgen dirdlg.go

type DirectoryDialog interface {
	Choose(initialDirectory string) (chosenDirectory string, err error)
}