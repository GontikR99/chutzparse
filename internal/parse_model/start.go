// +build wasm,electron

package parse_model

func StartMain() {
	listenForHits()
	listenForFights()
	maintainThroughput()
}