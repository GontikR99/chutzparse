// +build wasm,electron

package main

import "net/rpc"

var installers []func(server *rpc.Server)

func registerRpcHandler(installer func(server *rpc.Server)) {
	installers = append(installers, installer)
}

func newRpcServer() *rpc.Server {
	server := rpc.NewServer()
	for _, installer := range installers {
		installer(server)
	}
	return server
}
