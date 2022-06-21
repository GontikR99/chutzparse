//go:build wasm && electron
// +build wasm,electron

package model

func StartMain() {
	listenForHits()
	listenForFights()
	maintainThroughput()
}
