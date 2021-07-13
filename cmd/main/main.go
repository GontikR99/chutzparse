// +build electron,wasm

package main

import (
	"context"
	"github.com/gontikr99/chutzparse/internal"
	"github.com/gontikr99/chutzparse/internal/eqspec"
	"github.com/gontikr99/chutzparse/internal/iff"
	"github.com/gontikr99/chutzparse/internal/model"
	"github.com/gontikr99/chutzparse/internal/settings"
	"github.com/gontikr99/chutzparse/internal/ui"
	"github.com/gontikr99/chutzparse/pkg/console"
	"github.com/gontikr99/chutzparse/pkg/electron"
	"github.com/gontikr99/chutzparse/pkg/electron/application"
	"github.com/gontikr99/chutzparse/pkg/electron/browserwindow"
	"github.com/gontikr99/chutzparse/pkg/electron/screen"
	"github.com/gontikr99/chutzparse/pkg/nodejs/path"
	"time"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			console.Log(err)
			panic(err)
		}
	}()

	registerRpcHandler(eqspec.HandleRPC())
	registerRpcHandler(iff.HandleRPC())
	registerRpcHandler(settings.HandleRPC())
	registerRpcHandler(ui.HandleRPC())

	settings.DefaultSetting(settings.EverQuestDirectory, "C:\\Users\\Public\\Daybreak Game Company\\Installed Games\\EverQuest")

	appCtx, exitApp := context.WithCancel(context.Background())
	application.OnWindowAllClosed(exitApp)

	model.RegisterReports()
	model.StartMain()
	eqspec.RestartLogScans(appCtx)

	startup, ready := context.WithCancel(appCtx)
	application.OnReady(ready)
	<-startup.Done()

	mainBuilding, shown := context.WithCancel(appCtx)

	mainWindow := browserwindow.New(&browserwindow.Conf{
		Width:  1600,
		Height: 800,
		Show:   false,
		Title: "ChutzParse "+internal.Version,
		WebPreferences: &browserwindow.WebPreferences{
			Preload:          path.Join(application.GetAppPath(), "src/preload.js"),
			NodeIntegration:  false,
			ContextIsolation: true,
		},
	})
	mainWindow.OnClosed(exitApp)
	mainWindow.ServeRPC(newRpcServer())

	mainWindow.Once("ready-to-show", func() {
		mainWindow.RemoveMenu()
		mainWindow.Show()
		//mainWindow.JSValue().Get("webContents").Call("openDevTools", map[string]interface{} {
		//	"mode":"detach",
		//})
	})
	mainWindow.Once("show", shown)
	mainWindow.LoadFile(path.Join(application.GetAppPath(), "src/window.html"))
	<-mainBuilding.Done()

	primaryDisplay := screen.GetPrimaryDisplay()
	wndRect := electron.Rectangle{
		X:      primaryDisplay.Bounds.X + 100,
		Y:      primaryDisplay.Bounds.Y + 100,
		Width:  100,
		Height: 100,
	}

	overlayBuilding, oShown := context.WithCancel(appCtx)
	overlayWnd := browserwindow.New(&browserwindow.Conf{
		X:           wndRect.X,
		Y:           wndRect.Y,
		Title:       "ChutzParse Overlay",
		Width:       wndRect.Width,
		Height:      wndRect.Height,
		Show:        false,
		Transparent: true,
		Resizable:   false,
		Frame:       false,
		WebPreferences: &browserwindow.WebPreferences{
			Preload:          path.Join(application.GetAppPath(), "src/preload.js"),
			NodeIntegration:  false,
			ContextIsolation: true,
		},
	})
	overlayWnd.SetSkipTaskbar(true)
	overlayWnd.OnClosed(exitApp)
	overlayWnd.ServeRPC(newRpcServer())

	overlayWnd.Once("ready-to-show", func() {
		overlayWnd.RemoveMenu()
		overlayWnd.ShowInactive()
		overlayWnd.SetAlwaysOnTop(true)
		overlayWnd.SetIgnoreMouseEvents(true)
		//overlayWnd.JSValue().Get("webContents").Call("openDevTools", map[string]interface{} {
		//	"mode":"detach",
		//})
	})
	overlayWnd.Once("show", oShown)
	overlayWnd.LoadFile(path.Join(application.GetAppPath(), "src", "overlay.html"))
	<-overlayBuilding.Done()

	overlayShowing := true
	go func() {
		for {
			select {
			case <-appCtx.Done():
				return
			case <-time.After(50 * time.Millisecond):
				break
			}
			isTop := eqspec.IsTopWindow()
			if overlayShowing != isTop {
				if isTop {
					overlayShowing = true
					overlayWnd.ShowInactive()
				} else {
					overlayShowing = false
					overlayWnd.Hide()
					continue
				}
			}
			newLoc, err := eqspec.GetWindowExtents()
			if err != nil {
				continue
			}
			if wndRect != *newLoc {
				wndRect = *newLoc

				// A bit of a hack to account for high DPI support
				eqDisp := screen.GetDisplayMatching(&wndRect)

				overlayWnd.SetContentBounds(&electron.Rectangle{
					X:      primaryDisplay.Bounds.X + 100,
					Y:      primaryDisplay.Bounds.Y + 100,
					Width:  100,
					Height: 100,
				})
				<-time.After(50 * time.Millisecond)

				overlayWnd.SetContentBounds(&electron.Rectangle{
					X:      int(float64(wndRect.X) / primaryDisplay.ScaleFactor),
					Y:      int(float64(wndRect.Y) / primaryDisplay.ScaleFactor),
					Width:  int(float64(wndRect.Width) / eqDisp.ScaleFactor),
					Height: int(float64(wndRect.Height) / eqDisp.ScaleFactor),
				})
			}
		}
	}()

	<-appCtx.Done()
	console.Log("Exiting")
	browserwindow.CloseAll()
	application.Quit()

	<-context.Background().Done()
}
