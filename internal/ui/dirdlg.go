package ui

//go:generate ../../build/rpcgen dirdlg.go

// DirectoryDialog instructs the main process to run an "open directory" dialog.
type DirectoryDialog interface {
	Choose(initialDirectory string) (chosenDirectory string, err error)
}