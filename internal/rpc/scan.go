package rpc

import "net/rpc"

type RestartScanRequest struct{}
type RestartScanResponse struct{}

type StubRestartScan struct {
	restart func()
}

func (rss *StubRestartScan) RestartScan(req *RestartScanRequest, res *RestartScanResponse) error {
	rss.restart()
	return nil
}

func RestartScan(client *rpc.Client) error {
	return client.Call("StubRestartScan.RestartScan", new(RestartScanRequest), new(RestartScanResponse))
}

func HandleRestartScan(restart func()) func(server *rpc.Server) {
	return func(server *rpc.Server) {
		server.Register(&StubRestartScan{restart})
	}
}
