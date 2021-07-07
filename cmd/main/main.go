// +build electron,wasm

package main

import (
	"context"
	"github.com/gontikr99/chutzparse/cmd/main/mainrpc"
	"github.com/gontikr99/chutzparse/internal/eqlog"
	"github.com/gontikr99/chutzparse/internal/eqwnd"
	"github.com/gontikr99/chutzparse/internal/model"
	"github.com/gontikr99/chutzparse/internal/settings"
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

	settings.DefaultSetting(settings.EverQuestDirectory, "C:\\Users\\Public\\Daybreak Game Company\\Installed Games\\EverQuest")

	appCtx, exitApp := context.WithCancel(context.Background())
	application.OnWindowAllClosed(exitApp)

	model.RegisterReports()
	model.StartMain()
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
	overlayWnd.Once("show", oShown)
	overlayWnd.LoadFile(path.Join(application.GetAppPath(), "src", "overlay.html"))
	<-overlayBuilding.Done()

	go func() {
		for {
			select {
			case <-appCtx.Done():
				return
			case <-time.After(50 * time.Millisecond):
				break
			}
			newLoc, err := eqwnd.GetExtents()
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
	browserwindow.CloseAll()
	application.Quit()

	<-context.Background().Done()
}
