package ui

//go:generate ../../build/rpcgen clipboard.go

// Clipboard provides a means of instructing the main process to copy some text to the system clipboard
type Clipboard interface {
	Copy(text string) error
}
