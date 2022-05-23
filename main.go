package main

import (
	_ "embed"
	"github.com/getlantern/systray"
)

//go:embed assets/cam_white.png
var icon []byte

func onReady() {
	systray.SetIcon(icon)
	systray.SetTitle("")
	systray.SetTooltip("Screen video capture")
	mQuit := systray.AddMenuItem("Quit", "")
	mQuit.SetIcon(icon)

	go func() {
		for {
			select {
			case <-mQuit.ClickedCh:
				systray.Quit()
			}
		}
	}()
}

func onExit() {
	// clean up here
}

func main() {
	systray.Run(onReady, onExit)
	/*
		rec := capture.GetDisplayRect()
		capture.CaptureStart(func(pix uint8) {
			// write pixels
		}, capture.Options{rec.Width, rec.Height})
	*/
}
