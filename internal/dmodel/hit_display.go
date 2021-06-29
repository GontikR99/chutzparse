// +build wasm

package dmodel

type HitDisplayEvent struct {
	Text string
	Color string
	Big bool
}

const ChannelTopTarget="hitDisplayTop"
const ChannelBottomTarget="hitDisplayBottom"