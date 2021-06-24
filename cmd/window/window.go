// +build wasm

package main

import (
	"github.com/gontikr99/chutzparse/pkg/console"
	"github.com/gontikr99/chutzparse/pkg/vuguutil"
	"github.com/vugu/vugu"
	"github.com/vugu/vugu/domrender"
)

func main() {
	renderer, err := domrender.New("#page_root")
	if err != nil {
		panic(err)
	}
	defer renderer.Release()

	buildEnv, err := vugu.NewBuildEnv(renderer.EventEnv())
	if err != nil {
		panic(err)
	}

	root := &Root{}

	for ok := true; ok; ok = renderer.EventWait() {
		buildResults := buildEnv.RunBuild(root)

		err = renderer.Render(buildResults)
		vuguutil.InvokeRenderCallbacks()

		if err != nil {
			console.Log(err)
		}
	}
}
