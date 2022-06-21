//go:build wasm && web
// +build wasm,web

package ui

import "github.com/vugu/vugu"

type SelectBoxChangeEvent interface {
	Selected() map[string]struct{}
	//SetValue(string)
	Env() vugu.EventEnv
}

//vugugen:event SelectBoxChange
