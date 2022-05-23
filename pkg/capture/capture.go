package capture

import "errors"

type DisplayRect struct {
	Width, Height int
}

type Options struct {
	Width, Height int
}

type CallbackFrame func(pix uint8)

var isStarted bool = false

func (r *DisplayRect) AspectRationByWidth(w int) (width, height int) {
	return w, int(float32(r.Height) * (float32(w) / float32(r.Width)))
}

func CaptureStart(cb CallbackFrame, options Options) error {
	if isStarted {
		return errors.New("capture in progress now")
	}

	isStarted = true

	return nil
}

func StopCapture() {
	isStarted = false
}
