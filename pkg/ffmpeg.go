package pkg

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"time"
)

type RecordStatus uint

const (
	StatusStopped RecordStatus = iota
	StatusRecord
	StatusEmpty
)

type FFMpeg struct {
	cmd     *exec.Cmd
	chWrite chan []uint8
	buffer  []uint8
	stdin   io.WriteCloser
	status  chan RecordStatus
}

func NewFFMpeg() *FFMpeg {
	f := new(FFMpeg)
	f.status = make(chan RecordStatus, 1)
	return f
}

func (m *FFMpeg) Version() (string, error) {
	out, err := exec.Command("ffmpeg", "-version").Output()
	if err != nil {
		return "", err
	}
	re := regexp.MustCompile("ffmpeg version ([0-9.]+)")
	match := re.FindStringSubmatch(string(out))
	if len(match) >= 2 {
		return match[1], nil
	}
	return "", errors.New("can't get ffmpeg version")
}

func (m *FFMpeg) statusEmpty() {
	m.status <- StatusEmpty
}

func (m *FFMpeg) statusRecord() {
	m.status <- StatusRecord
}

func (m *FFMpeg) statusStopped() {
	m.status <- StatusStopped
}

func (m *FFMpeg) handleStatus() {
	for {
		select {
		case status := <-m.status:
			switch status {
			case StatusRecord:
				// onStart
			case StatusStopped:
				// Wait for frames
				for len(m.chWrite) > 0 {
					time.Sleep(time.Second)
				}
				// Close output
				m.stdin.Close()
			case StatusEmpty:
				// onStop
				return
			}
		}
	}
}

func (m *FFMpeg) Write(pix []uint8) {
	m.chWrite <- pix
}

func (m *FFMpeg) Record(width, height int, out string) error {
	m.cmd = exec.Command("ffmpeg",
		"-f", "rawvideo",
		"-framerate", "60",
		"-pix_fmt", "bgra",
		"-video_size", fmt.Sprintf("%vx%v", width, height),
		"-i", "-",
		"-vcodec", "mpeg4",
		"-q:v", "0",
		"-quality", "best",
		"-r", "60",
		"-crf", "0",
		out,
	)

	var err error
	m.stdin, err = m.cmd.StdinPipe()
	if err != nil {
		return err
	}
	m.chWrite = make(chan []uint8, 1)

	go func() {
		out, _ := m.cmd.CombinedOutput()
		fmt.Printf("FFMPEG output:\n%v\n", string(out))
		m.statusEmpty()
	}()

	go m.handleStatus()
	go m.handleAsyncWriter(width, height)
	go m.handleFrames(width, height)

	m.statusRecord()

	return nil
}

func (m *FFMpeg) handleAsyncWriter(width, height int) {
	buf := new(bytes.Buffer)
	tick := time.Tick(16 * time.Millisecond)
	for {
		select {
		case <-tick:
			stride := len(m.buffer) / height
			rowLen := 4 * width
			for i := 0; i < len(m.buffer); i += stride {
				if _, err := m.stdin.Write(m.buffer[i : i+rowLen]); err != nil {
					break
				}
			}
			buf.Reset()
			// fmt.Print(".")
		}
	}
}

func (m *FFMpeg) handleFrames(width, height int) {
	// buf := new(bytes.Buffer)
	for {
		select {
		case pix := <-m.chWrite:
			m.buffer = pix
			//stride := len(pix) / height
			//rowLen := 4 * width
			//for i := 0; i < len(pix); i += stride {
			//	if _, err := m.stdin.Write(pix[i : i+rowLen]); err != nil {
			//		break
			//	}
			//}
			// buf.Reset()
		}
	}
}

func (m *FFMpeg) Stop() {
	m.statusStopped()
}
