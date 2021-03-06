// +build wasm

package ipc

import "net/rpc"

const ChannelRPCMain = "rpcMain"

// NewClient creates a new RPC client on the specified endpoint/channel name
func NewClient(channelName string, endpoint Endpoint) *rpc.Client {
	return rpc.NewClient(EndpointAsStream(channelName, endpoint))
}
