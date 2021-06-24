// +build wasm,web

package main

import (
	"github.com/gontikr99/chutzparse/cmd/overlay/throughput"
	"github.com/gontikr99/chutzparse/pkg/console"
	"github.com/vugu/vugu"
)

type Root struct {
}


func (c *Root) Init(ctx vugu.InitCtx) {
	console.Log("Init root")
}

var topBars = []throughput.Bar{
	{"oh", "Joram", "god", []throughput.BarSector{{"blue", 0.8}, {"cyan", 0.2}}},
	{"my", "Gontik", "heals!", []throughput.BarSector{{"blue", 0.6}, {"cyan", 0.25}}},
}

var bottomBars = []throughput.Bar{
	{"", "Siluine", "100%", []throughput.BarSector{{"blue", 1}}},
	{"Me", "Jephine", "80%", []throughput.BarSector{{"red", 0.8}}},
	{"", "Joramini", "40%", []throughput.BarSector{{"blue", 0.4}}},
}