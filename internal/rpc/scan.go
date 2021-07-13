package rpc

//go:generate ../../build/rpcgen scan.go

type ScanControl interface {
	Restart() error
}
