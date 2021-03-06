package pkg

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strconv"
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
	fps     int
	chWrite chan []uint8
	buffer  []uint8
	stdin   io.WriteCloser
	status  chan RecordStatus
	work    bool
}

func NewFFMpeg() *FFMpeg {
	f := new(FFMpeg)
	f.status = make(chan RecordStatus, 1)
	f.work = false
	f.fps = 60
	return f
}

func (m *FFMpeg) SetFPS(f int) {
	m.fps = f
}

func (m *FFMpeg) Version() (string, error) {
	out, err := exec.Command("ffmpeg", "-version").Output()
	if err != nil {
		return "", err
	}
	re := regexp.MustCompile("ffmpeg version ([A-z0-9.-]+)")
	match := re.FindStringSubmatch(string(out))
	fmt.Println(match)
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
				m.work = true
			case StatusStopped:
				for len(m.chWrite) > 0 {
					time.Sleep(time.Second)
				}
				m.stdin.Close()
			case StatusEmpty:
				m.work = false
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
		"-framerate", strconv.Itoa(m.fps),
		"-pix_fmt", "bgra",
		"-video_size", fmt.Sprintf("%vx%v", width, height),
		"-i", "-",
		"-vcodec", "mpeg4",
		"-q:v", "1",
		"-r", strconv.Itoa(m.fps),
		// "-crf", "22",
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
	go m.handleFrames()

	m.statusRecord()

	return nil
}

func (m *FFMpeg) handleAsyncWriter(width, height int) {
	buf := new(bytes.Buffer)
	tick := time.Tick(time.Duration(1000/m.fps) * time.Millisecond)
	for {
		if m.work == false {
			return
		}
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
		}
	}
}

func (m *FFMpeg) handleFrames() {
	for {
		if m.work == false {
			return
		}
		select {
		case pix := <-m.chWrite:
			m.buffer = pix
		}
	}
}

func (m *FFMpeg) Stop() {
	m.statusStopped()
}
