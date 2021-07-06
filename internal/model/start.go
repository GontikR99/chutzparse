// +build wasm,electron

package model

func StartMain() {
	listenForHits()
	listenForFights()
	maintainThroughput()
}
