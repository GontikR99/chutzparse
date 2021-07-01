// +build wasm,electron

package screen

import "github.com/gontikr99/chutzparse/pkg/electron"

var scrJs = electron.JSValue().Get("screen")

func GetPrimaryDisplay() *electron.Display {
	return electron.JSValueToDisplay(scrJs.Call("getPrimaryDisplay"))
}

func GetDisplayNearestPoint(p *electron.Point) *electron.Display {
	return electron.JSValueToDisplay(scrJs.Call("getDisplayNearestPoint", p.JSValue()))
}

func ScreenToDipPoint(p *electron.Point) *electron.Point {
	return electron.JSValueToPoint(scrJs.Call("screenToDipPoint", p.JSValue()))
}

func DipToScreenPoint(p *electron.Point) *electron.Point {
	return electron.JSValueToPoint(scrJs.Call("dipToScreenPoint", p.JSValue()))
}

func GetDisplayMatching(r *electron.Rectangle) *electron.Display {
	return electron.JSValueToDisplay(scrJs.Call("getDisplayMatching", r.JSValue()))
}