// +build wasm,web

package main

import (
	"github.com/gontikr99/chutzparse/pkg/console"
	"github.com/vugu/vugu"
)

type Root struct {
}


func (c *Root) Init(ctx vugu.InitCtx) {
	console.Log("Init root")
}
