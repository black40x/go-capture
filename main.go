package main

import "github.com/black40x/go-capture/pkg/capture"

func main() {
	rec := capture.GetDisplayRect()
	capture.CaptureStart(func(pix uint8) {
		// write pixels
	}, capture.Options{rec.Width, rec.Height})
}
