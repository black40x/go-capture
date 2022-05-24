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

type Application struct {
	ffmpeg  *FFMpeg
	display *capture.DisplayRect
}

func NewApplication() *Application {
	app := new(Application)
	app.ffmpeg = NewFFMpeg()
	app.display = capture.GetDisplayRect()
	return app
}

func (a *Application) trayOnReady() {
	systray.SetIcon(icon)
	systray.SetTitle("")
	systray.SetTooltip("Screen video capture")
	//
	ffVer, ferr := a.ffmpeg.Version()
	mAbout := systray.AddMenuItem("About", "")
	mFFmpeg := systray.AddMenuItem("Install ffmpeg", "")
	if ferr == nil {
		mFFmpeg.SetTitle(fmt.Sprintf("FFmpeg v%s", ffVer))
	}

	systray.AddSeparator()
	mCapture := systray.AddMenuItem("Capture", "")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "")
	//

	// Actions
	go func() {
		for {
			select {
			case <-mAbout.ClickedCh:
				fmt.Println("About")
			case <-mQuit.ClickedCh:
				systray.Quit()
			case <-mCapture.ClickedCh:
				if !capture.IsStarted() {
					mCapture.SetTitle("Stop")
					a.captureStart()
				} else {
					mCapture.SetTitle("Capture")
					a.captureStop()
				}
			case <-mFFmpeg.ClickedCh:
				if ferr != nil {
					fmt.Println("GoTo install ffmpeg")
				} else {
					fmt.Println("GoTo ffmpeg")
				}
			}
		}
	}()
}

func (a *Application) trayOnExit() {
	//
}

func (a *Application) captureStop() {
	capture.CaptureStop()
	a.ffmpeg.Stop()
}

var fps int
var tstart time.Time

func (a *Application) captureStart() {
	tstart = time.Now()
	fps = 0

	a.ffmpeg.Record(a.display.Width, a.display.Height, "out.mov")
	capture.CaptureStart(func(pix []uint8, time uint64) {
		fps = 1000 / int(time/1000000)
		a.ffmpeg.Write(pix)
	}, capture.Options{Width: a.display.Width, Height: a.display.Height})

	ticker := time.NewTicker(time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				systray.SetTitle(fmt.Sprintf("%d FPS | %s", fps, fmtDuration(time.Since(tstart))))
			}
		}
	}()
	// ToDO : chan and errors
}

func fmtDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= h * time.Minute
	s := d / time.Second
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

func (a *Application) Exec() {
	systray.Run(a.trayOnReady, a.trayOnExit)
}
