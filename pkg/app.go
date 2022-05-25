package pkg

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/black40x/go-capture/pkg/capture"
	"github.com/getlantern/systray"
)

//go:embed assets/cam_white.png
var icon []byte

//go:embed assets/cam_white_rec.png
var iconRecord []byte

//go:embed assets/cam_white_rec.ico
var iconIco []byte

//go:embed assets/cam_white_rec.ico
var iconRecordIco []byte

type Application struct {
	ffmpeg  *FFMpeg
	display *capture.DisplayRect
	// Menus
	mAbout, mCapture, mQuit               *systray.MenuItem
	mFFMpeg, mFFMpegInstall, mFFMpegAbout *systray.MenuItem
	mFPS, m60FPS, m30FPS                  *systray.MenuItem
}

func NewApplication() *Application {
	app := new(Application)
	app.ffmpeg = NewFFMpeg()
	app.display = capture.GetDisplayRect()
	return app
}

func (a *Application) openURL(url string) {
	switch runtime.GOOS {
	case "linux":
		exec.Command("xdg-open", url).Start()
	case "windows":
		exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		exec.Command("open", url).Start()
	}
}

func (a *Application) buildMenu() {
	ffVer, ferr := a.ffmpeg.Version()

	a.mAbout = systray.AddMenuItem("About", "")
	a.mFFMpeg = systray.AddMenuItem("FFmpeg", "")
	a.mFFMpegInstall = systray.AddMenuItem("FFmpeg install", "")
	a.mFFMpeg.SetTitle(fmt.Sprintf("FFmpeg v%s", ffVer))

	a.mFPS = a.mFFMpeg.AddSubMenuItem("FPS", "")
	a.m60FPS = a.mFPS.AddSubMenuItemCheckbox("60 FPS", "", true)
	a.m30FPS = a.mFPS.AddSubMenuItemCheckbox("30 FPS", "", false)
	systray.AddSeparator()
	a.mCapture = systray.AddMenuItem("Capture", "")
	a.mFFMpegAbout = a.mFFMpeg.AddSubMenuItem("About", "")

	if ferr == nil {
		a.mFFMpeg.Show()
		a.mFFMpegInstall.Hide()
	} else {
		a.mCapture.Hide()
		a.mFFMpeg.Hide()
		a.mFFMpegInstall.Show()
	}
	systray.AddSeparator()
	a.mQuit = systray.AddMenuItem("Quit", "")
}

func (a *Application) handleMenuActions() {
	for {
		select {
		case <-a.mFFMpegInstall.ClickedCh:
			a.openURL("https://ffmpeg.org/")
		case <-a.mFFMpegAbout.ClickedCh:
			a.openURL("https://ffmpeg.org/")
		case <-a.mAbout.ClickedCh:
			a.openURL("https://github.com/black40x/go-capture")
		case <-a.mQuit.ClickedCh:
			systray.Quit()
		//
		case <-a.m60FPS.ClickedCh:
			a.ffmpeg.SetFPS(60)
			a.m60FPS.Check()
			a.m30FPS.Uncheck()
		case <-a.m30FPS.ClickedCh:
			a.ffmpeg.SetFPS(30)
			a.m30FPS.Check()
			a.m60FPS.Uncheck()
		//
		case <-a.mCapture.ClickedCh:
			if !capture.IsStarted() {
				if a.captureStart() != nil {
				} else {
					a.mCapture.SetTitle("Stop")
					a.mFPS.Disable()
				}
			} else {
				a.mCapture.SetTitle("Capture")
				a.mFPS.Enable()
				a.captureStop()
			}
		}
	}
}

func (a *Application) trayOnReady() {
	if runtime.GOOS == "windows" {
		systray.SetIcon(iconIco)
	} else {
		systray.SetIcon(icon)
	}
	systray.SetTitle("")
	systray.SetTooltip("Screen video capture")
	a.buildMenu()
	go a.handleMenuActions()
}

func (a *Application) trayOnExit() {
	//
}

func (a *Application) captureStop() {
	if runtime.GOOS == "windows" {
		systray.SetIcon(iconIco)
	} else {
		systray.SetIcon(icon)
	}
	capture.CaptureStop()
	a.ffmpeg.Stop()
}

func (a *Application) captureStart() error {
	c, _ := os.UserHomeDir()
	fn := fmt.Sprintf("%s/Desktop/ScreenCapture_%s.mov", c, time.Now().Format("01_01_2006_15_04_05"))
	fmt.Println(fn)
	err := a.ffmpeg.Record(a.display.Width, a.display.Height, fn)
	if err != nil {
		return err
	}

	err = capture.CaptureStart(func(pix []uint8, time uint64) {
		a.ffmpeg.Write(pix)
	}, capture.Options{Width: a.display.Width, Height: a.display.Height})
	if err != nil {
		a.ffmpeg.Stop()
		return err
	}

	if runtime.GOOS == "windows" {
		systray.SetIcon(iconRecordIco)
	} else {
		systray.SetIcon(iconRecord)
	}

	return nil
}

func (a *Application) Exec() {
	systray.Run(a.trayOnReady, a.trayOnExit)
}
