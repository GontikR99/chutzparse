package electron

import "syscall/js"

type Point struct {
	X int
	Y int
}

func (p *Point) JSValue() js.Value {
	return js.ValueOf(map[string]interface{}{
		"x":p.X,
		"y":p.Y,
	})
}

func JSValueToPoint(value js.Value) *Point {
	return &Point{
		X: value.Get("x").Int(),
		Y: value.Get("y").Int(),
	}
}

type Size struct {
	Width int
	Height int
}

func (s *Size) JSValue() js.Value {
	return js.ValueOf(map[string]interface{}{
		"width":s.Width,
		"height":s.Height,
	})
}

func JSValueToSize(value js.Value) *Size {
	return &Size{
		Width: value.Get("width").Int(),
		Height: value.Get("height").Int(),
	}
}

type Rectangle struct {
	X      int
	Y      int
	Width  int
	Height int
}

func (rectangle *Rectangle) JSValue() js.Value {
	return js.ValueOf(map[string]interface{}{
		"x":      rectangle.X,
		"y":      rectangle.Y,
		"width":  rectangle.Width,
		"height": rectangle.Height,
	})
}

func JSValueToRectangle(value js.Value) *Rectangle {
	return &Rectangle{
		X:      value.Get("x").Int(),
		Y:      value.Get("y").Int(),
		Width:  value.Get("width").Int(),
		Height: value.Get("height").Int(),
	}
}

type Display struct {
	Id int
	Rotation int
	ScaleFactor float64
	TouchSupport string
	Monochrome bool
	AccelerometerSupport string
	ColorSpace string
	ColorDepth int
	DepthPerComponent int
	DisplayFrequency int
	Bounds *Rectangle
	Size *Size
	WorkArea *Rectangle
	WorkAreaSize *Size
	Internal bool
}

func JSValueToDisplay(value js.Value) *Display {
	return &Display{
		Id:                   value.Get("id").Int(),
		Rotation:             value.Get("rotation").Int(),
		ScaleFactor:          value.Get("scaleFactor").Float(),
		TouchSupport:         value.Get("touchSupport").String(),
		Monochrome:           value.Get("monochrome").Bool(),
		AccelerometerSupport: value.Get("accelerometerSupport").String(),
		ColorSpace:           value.Get("colorSpace").String(),
		ColorDepth:           value.Get("colorDepth").Int(),
		DepthPerComponent:    value.Get("depthPerComponent").Int(),
		DisplayFrequency:     value.Get("displayFrequency").Int(),
		Bounds:               JSValueToRectangle(value.Get("bounds")),
		Size:                 JSValueToSize(value.Get("size")),
		WorkArea:             JSValueToRectangle(value.Get("workArea")),
		WorkAreaSize:         JSValueToSize(value.Get("workAreaSize")),
		Internal:             value.Get("internal").Bool(),
	}
}