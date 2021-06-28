// +build wasm

package dmodel

type HitDisplayEvent struct {
	Text string
	Color string
}

const ChannelTopTarget="hitDisplayTop"
const ChannelBottomTarget="hitDisplayBottom"