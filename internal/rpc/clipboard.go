package rpc

//go:generate ../../build/rpcgen clipboard.go

type Clipboard interface {
	Copy(text string) error
}