package eqspec

//go:generate ../../build/rpcgen scancontrol.go

// ScanControl instructs the main process to restart log scanning
type ScanControl interface {
	Restart() error
}

