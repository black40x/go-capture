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
	//
	//mCapture := systray.AddMenuItem("Capture", "")
	systray.AddSeparator()
	//mFFmpeg := systray.AddMenuItem("FFmpeg", "")
	//mFFmpeg.AddSubMenuItemCheckbox()
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "")

	// Actions
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
}
