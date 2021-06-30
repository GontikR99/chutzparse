// +build electron,wasm

package main

import (
	"context"
	"github.com/gontikr99/chutzparse/cmd/main/mainrpc"
	"github.com/gontikr99/chutzparse/internal/eqlog"
	"github.com/gontikr99/chutzparse/internal/eqwnd"
	"github.com/gontikr99/chutzparse/internal/parse_model"
	"github.com/gontikr99/chutzparse/internal/settings"
	"github.com/gontikr99/chutzparse/pkg/console"
	"github.com/gontikr99/chutzparse/pkg/electron"
	"github.com/gontikr99/chutzparse/pkg/electron/application"
	"github.com/gontikr99/chutzparse/pkg/electron/browserwindow"
	"github.com/gontikr99/chutzparse/pkg/nodejs/path"
	"time"
)

func main() {
	defer func(){
		if err := recover(); err!=nil {
			console.Log(err)
			panic(err)
		}
	}()

	application.JSValue().Get("commandLine").Call("appendSwitch", "high-dpi-support", 1)
	application.JSValue().Get("commandLine").Call("appendSwitch", "force-device-scale-factor", 1)

	settings.DefaultSetting(settings.EverQuestDirectory, "C:\\Users\\Public\\Daybreak Game Company\\Installed Games\\EverQuest")

	appCtx, exitApp := context.WithCancel(context.Background())
	application.OnWindowAllClosed(exitApp)

	parse_model.StartMain()
	eqlog.RestartLogScans(appCtx)

	startup, ready := context.WithCancel(appCtx)
	application.OnReady(ready)
	<-startup.Done()

	mainBuilding, shown := context.WithCancel(appCtx)

	mainWindow := browserwindow.New(&browserwindow.Conf{
		Width:  1600,
		Height: 800,
		Show:   false,
		WebPreferences: &browserwindow.WebPreferences{
			Preload:          path.Join(application.GetAppPath(), "src/preload.js"),
			NodeIntegration:  false,
			ContextIsolation: true,
		},
	})
	mainWindow.OnClosed(exitApp)
	mainWindow.ServeRPC(mainrpc.NewServer())

	mainWindow.Once("ready-to-show", func() {
		mainWindow.RemoveMenu()
		mainWindow.Show()
		shown()
	})
	mainWindow.LoadFile(path.Join(application.GetAppPath(), "src/window.html"))
	<- mainBuilding.Done()

	wndRect := electron.Rectangle{
		X:      100,
		Y:      100,
		Width:  100,
		Height: 100,
	}
	overlayWnd := browserwindow.New(&browserwindow.Conf{
		X:              wndRect.X,
		Y:              wndRect.Y,
		Title:          "ChutzParse Overlay",
		Width:          wndRect.Width,
		Height:         wndRect.Height,
		Show:           false,
		Transparent:    true,
		Resizable:      false,
		Frame:          false,
		WebPreferences: &browserwindow.WebPreferences{
			Preload:          path.Join(application.GetAppPath(), "src/preload.js"),
			NodeIntegration:  false,
			ContextIsolation: true,
		},
	})
	overlayWnd.OnClosed(exitApp)
	overlayWnd.ServeRPC(mainrpc.NewServer())

	overlayWnd.Once("ready-to-show", func() {
		overlayWnd.RemoveMenu()
		overlayWnd.ShowInactive()
		overlayWnd.SetAlwaysOnTop(true)
		overlayWnd.SetIgnoreMouseEvents(true)
		//overlayWnd.JSValue().Get("webContents").Call("openDevTools", map[string]interface{} {
		//	"mode":"detach",
		//})
	})
	overlayWnd.LoadFile(path.Join(application.GetAppPath(), "src","overlay.html"))

	go func() {
		for {
			select {
			case <-appCtx.Done():
				return
			case <-time.After(50*time.Millisecond):
				break
			}
			newLoc, err := eqwnd.GetExtents()
			if err!=nil {
				continue
			}
			if wndRect != *newLoc {
				wndRect = *newLoc
				overlayWnd.SetContentBounds(&wndRect)
			}
		}
	}()

	<-appCtx.Done()
	browserwindow.CloseAll()
	application.Quit()

	<-context.Background().Done()
}