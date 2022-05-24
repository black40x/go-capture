package pkg

import (
	_ "embed"
	"fmt"
	"github.com/black40x/go-capture/pkg/capture"
	"github.com/getlantern/systray"
	"time"
)

//go:embed assets/cam_white.png
var icon []byte

//go:embed assets/cam_white_record.png
var iconRecord []byte

type Application struct {
	ffmpeg  *FFMpeg
	display *capture.DisplayRect
	//
	fps int
	// Menus
	mAbout, mCapture, mQuit               *systray.MenuItem
	mFFMpeg, mFFMpegInstall, mFFMpegAbout *systray.MenuItem
	mFPS, m60FPS, m30FPS                  *systray.MenuItem
}

func NewApplication() *Application {
	app := new(Application)
	app.ffmpeg = NewFFMpeg()
	app.display = capture.GetDisplayRect()
	app.fps = 60

	return app
}

func (a *Application) buildMenu() {
	ffVer, ferr := a.ffmpeg.Version()

	a.mAbout = systray.AddMenuItem("About", "")
	a.mFFMpeg = systray.AddMenuItem("FFmpeg", "")
	a.mFFMpegInstall = systray.AddMenuItem("FFmpeg install", "")

	if ferr == nil {
		a.mFFMpeg.SetTitle(fmt.Sprintf("FFmpeg v%s", ffVer))

		a.mFPS = a.mFFMpeg.AddSubMenuItem("FPS", "")
		a.m60FPS = a.mFPS.AddSubMenuItemCheckbox("60 FPS", "", true)
		a.m30FPS = a.mFPS.AddSubMenuItemCheckbox("30 FPS", "", false)

		a.mFFMpegAbout = a.mFFMpeg.AddSubMenuItem("About", "")

		a.mFFMpeg.Show()
		a.mFFMpegInstall.Hide()
		systray.AddSeparator()
		a.mCapture = systray.AddMenuItem("Capture", "")
	} else {
		a.mFFMpeg.Hide()
		a.mFFMpegInstall.Show()
	}
	systray.AddSeparator()
	a.mQuit = systray.AddMenuItem("Quit", "")
}

func (a *Application) handleMenuActions() {
	for {
		select {
		case <-a.mFFMpegAbout.ClickedCh:
			fmt.Println("About ffmpeg")
		case <-a.mAbout.ClickedCh:
			fmt.Println("About")
		case <-a.mQuit.ClickedCh:
			systray.Quit()
		//
		case <-a.m60FPS.ClickedCh:
			a.fps = 60
			a.m60FPS.Check()
			a.m30FPS.Uncheck()
		case <-a.m30FPS.ClickedCh:
			a.fps = 30
			a.m30FPS.Check()
			a.m60FPS.Uncheck()
		//
		case <-a.mCapture.ClickedCh:
			if !capture.IsStarted() {
				if a.captureStart() != nil {
				} else {
					a.mCapture.SetTitle("Stop")
				}
			} else {
				a.mCapture.SetTitle("Capture")
				a.captureStop()
			}
		case <-a.mFFMpegInstall.ClickedCh:
			fmt.Println("GoTo ffmpeg")
		}
	}
}

func (a *Application) trayOnReady() {
	systray.SetIcon(icon)
	systray.SetTitle("")
	systray.SetTooltip("Screen video capture")
	a.buildMenu()
	go a.handleMenuActions()
}

func (a *Application) trayOnExit() {
	//
}

func (a *Application) captureStop() {
	systray.SetIcon(icon)
	capture.CaptureStop()
	a.ffmpeg.Stop()
}

func (a *Application) captureStart() error {
	fn := fmt.Sprintf("ScreenCapture_%s.mov", time.Now().Format("01_01_2006_15_04_05"))
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

	systray.SetIcon(iconRecord)

	return nil
}

func (a *Application) Exec() {
	systray.Run(a.trayOnReady, a.trayOnExit)
}
