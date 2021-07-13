package rpc

//go:generate ../../build/rpcgen scan.go

// ScanControl instructs the main process to restart log scanning
type ScanControl interface {
	Restart() error
}
