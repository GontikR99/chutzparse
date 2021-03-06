// +build wasm,electron

package browserwindow

import (
	"github.com/gontikr99/chutzparse/pkg/console"
	"github.com/gontikr99/chutzparse/pkg/electron"
	"github.com/gontikr99/chutzparse/pkg/electron/binding"
	"github.com/gontikr99/chutzparse/pkg/electron/ipc"
	"github.com/gontikr99/chutzparse/pkg/electron/ipc/ipcmain"
	"github.com/gontikr99/chutzparse/pkg/jsbinding"
	"io"
	"net/rpc"
	"strconv"
	"syscall/js"
)

var browserWindow = electron.JSValue().Get("BrowserWindow")

type BrowserWindow interface {
	io.Closer
	ipc.Endpoint

	Destroy()

	RemoveMenu()
	Show()
	ShowInactive()
	Hide()
	LoadFile(path string)

	Id() int

	On(eventName string, action func())
	Once(eventName string, action func())

	SetAlwaysOnTop(bool)

	GetBounds() *electron.Rectangle
	SetBounds(rectangle *electron.Rectangle)

	GetContentBounds() *electron.Rectangle
	SetContentBounds(rectangle *electron.Rectangle)

	SetIgnoreMouseEvents(bool)
	SetSkipTaskbar(bool)

	// Additions
	OnClosed(action func())
	ServeRPC(server *rpc.Server)
	JSValue() js.Value
}

type electronBrowserWindow struct {
	browserWindow js.Value
	webContents   js.Value

	windowId int

	closedCallbacks []func()
	callbacks       map[int]js.Func
	nextCallback    int
}

var nextWindowId = 0
var openWindows = make(map[int]*electronBrowserWindow)

type WebPreferences struct {
	Preload          interface{} `json:"preload"`
	ContextIsolation interface{} `json:"contextIsolation"`
	NodeIntegration  interface{} `json:"nodeIntegration"`
}

type Conf struct {
	X              interface{}     `json:"x"`
	Y              interface{}     `json:"y"`
	Title          interface{}     `json:"title"`
	Width          interface{}     `json:"width"`
	Height         interface{}     `json:"height"`
	Show           interface{}     `json:"show"`
	Transparent    interface{}     `json:"transparent"`
	Resizable      interface{}     `json:"resizable"`
	Frame          interface{}     `json:"frame"`
	WebPreferences *WebPreferences `json:"webPreferences"`
}

func New(conf *Conf) BrowserWindow {
	if nextWindowId < 0 {
		console.Log("All windows should be closed")
		return nil
	}
	browserWindowInternal := browserWindow.New(binding.JsonifyOptions(conf))

	browserWindowInstance := &electronBrowserWindow{
		browserWindow: browserWindowInternal,
		webContents:   browserWindowInternal.Get("webContents"),
		callbacks:     make(map[int]js.Func),
		nextCallback:  0,
		windowId:      nextWindowId,
	}
	nextWindowId++
	openWindows[nextWindowId] = browserWindowInstance

	var handleClosedFunc js.Func
	handleClosedFuncAddr := &handleClosedFunc
	*handleClosedFuncAddr = js.FuncOf(func(_ js.Value, _ []js.Value) interface{} {
		browserWindowInstance.handleClosed()
		(*handleClosedFuncAddr).Release()
		return nil
	})
	browserWindowInternal.Call("on", "closed", handleClosedFunc)

	return browserWindowInstance
}

func CloseAll() {
	nextWindowId = -1
	for _, v := range openWindows {
		v.Destroy()
	}
}

func (bw *electronBrowserWindow) registerCallback(callback func()) (int, js.Func) {
	wrapped := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		callback()
		return nil
	})
	bw.nextCallback++
	bw.callbacks[bw.nextCallback] = wrapped
	return bw.nextCallback, wrapped
}

func (bw *electronBrowserWindow) singleshotCallback(callback func()) js.Func {
	var cbHolder int
	cbLoc := &cbHolder
	cbId, wrapped := bw.registerCallback(func() {
		callback()
		if fnc, ok := bw.callbacks[*cbLoc]; ok {
			fnc.Release()
			delete(bw.callbacks, *cbLoc)
		}
	})
	*cbLoc = cbId
	return wrapped
}

// On window closed, call all of our queued close callbacks, then release any outstanding callbacks
func (bw *electronBrowserWindow) handleClosed() {
	for i := len(bw.closedCallbacks) - 1; i >= 0; i-- {
		bw.closedCallbacks[i]()
	}
	for _, w := range bw.callbacks {
		w.Release()
	}
	delete(openWindows, bw.windowId)
}

// Do something when the window is closed
func (bw *electronBrowserWindow) OnClosed(callback func()) {
	bw.closedCallbacks = append(bw.closedCallbacks, callback)
}

// Close the window, as if the user had clicked the close button.  May be intercepted/interrupted by code
func (bw *electronBrowserWindow) Close() error {
	bw.browserWindow.Call("close")
	return nil
}

// Close the window with no chance for interception.
func (bw *electronBrowserWindow) Destroy() {
	bw.browserWindow.Call("destroy")
}

// Turn off the menu bar
func (bw *electronBrowserWindow) RemoveMenu() {
	bw.browserWindow.Call("removeMenu")
}

// Show the window
func (bw *electronBrowserWindow) Show() {
	bw.browserWindow.Call("show")
}

// Show the window
func (bw *electronBrowserWindow) ShowInactive() {
	bw.browserWindow.Call("showInactive")
}

// Load content into the window
func (bw *electronBrowserWindow) LoadFile(path string) {
	bw.browserWindow.Call("loadFile", path)
}

// Register a callback to be called once
func (bw *electronBrowserWindow) Once(eventName string, action func()) {
	if eventName == "closed" {
		bw.closedCallbacks = append(bw.closedCallbacks, action)
	} else {
		bw.browserWindow.Call("once", eventName, bw.singleshotCallback(action))
	}
}

// Register a callback to be called repeatedly
func (bw *electronBrowserWindow) On(eventName string, action func()) {
	if eventName == "closed" {
		bw.closedCallbacks = append(bw.closedCallbacks, action)
	} else {
		_, wrapped := bw.registerCallback(action)
		bw.browserWindow.Call("on", eventName, wrapped)
	}
}

func (bw *electronBrowserWindow) SetAlwaysOnTop(b bool) {
	bw.browserWindow.Call("setAlwaysOnTop", b)
}

// Send a message to this window on the specified channel
func (bw *electronBrowserWindow) Send(channelName string, content []byte) {
	defer func() {
		if r := recover(); r != nil {
			// ignore errors sending
		}
	}()
	bw.webContents.Call("send", ipc.Prefix+channelName, jsbinding.MakeArrayBuffer(content))
}

// Send a message to all open windows on the specified channel
func Broadcast(channel string, content []byte) {
	for _, v := range openWindows {
		v.Send(channel, content)
	}
}

func BroadcastChunked(channel string, content []byte) {
	for _, v := range openWindows {
		ipc.SendChunked(v, channel, content)
	}
}

func (bw *electronBrowserWindow) Id() int {
	return bw.webContents.Get("id").Int()
}

func (bw *electronBrowserWindow) Listen(channelName string) (<-chan ipc.Message, func()) {
	outChan := make(chan ipc.Message)
	inChan, inDone := ipcmain.Listen(channelName)
	go func() {
		for {
			select {
			case inMsg := <-inChan:
				if inMsg == nil {
					return
				}
				if inMsg.Sender() == strconv.Itoa(bw.Id()) {
					outChan <- inMsg
				}
			}
		}
	}()
	return outChan, inDone
}

func (bw *electronBrowserWindow) ServeRPC(server *rpc.Server) {
	endpointStream := ipc.EndpointAsStream("rpcMain", bw)
	bw.OnClosed(func() {
		endpointStream.Close()
	})
	go func() {
		server.ServeConn(endpointStream)
	}()
}

func (bw *electronBrowserWindow) JSValue() js.Value {
	return bw.browserWindow
}

func (bw *electronBrowserWindow) GetBounds() *electron.Rectangle {
	return electron.JSValueToRectangle(bw.browserWindow.Call("getBounds"))
}

func (bw *electronBrowserWindow) SetBounds(rectangle *electron.Rectangle) {
	bw.browserWindow.Call("setBounds", rectangle.JSValue())
}

func (bw *electronBrowserWindow) GetContentBounds() *electron.Rectangle {
	return electron.JSValueToRectangle(bw.browserWindow.Call("getContentBounds"))
}

func (bw *electronBrowserWindow) SetContentBounds(rectangle *electron.Rectangle) {
	bw.browserWindow.Call("setContentBounds", rectangle.JSValue())
}

func (bw *electronBrowserWindow) SetIgnoreMouseEvents(b bool) {
	bw.browserWindow.Call("setIgnoreMouseEvents", b)
}

func (bw *electronBrowserWindow) SetSkipTaskbar(b bool) {
	bw.browserWindow.Call("setSkipTaskbar", b)
}

func (bw *electronBrowserWindow) Hide() {
	bw.browserWindow.Call("hide")
}
