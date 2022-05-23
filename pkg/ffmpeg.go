package pkg

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"regexp"
)

type FFMpeg struct {
	cmd     *exec.Cmd
	chWrite chan []uint8
	stdin   io.WriteCloser
	stop    bool
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

func (m *FFMpeg) Write(pix []uint8) {
	m.chWrite <- pix
}

func (m *FFMpeg) Record(width, height int, out string) error {
	m.cmd = exec.Command("ffmpeg",
		"-f", "rawvideo",
		"-pix_fmt", "bgra",
		"-video_size", fmt.Sprintf("%vx%v", width, height),
		"-i", "-",
		"-vcodec", "mpeg4",
		"-q:v", "0",
		"-quality", "best",
		"-r", "30",
		"-crf", "0",
		out,
	)

	var err error
	m.stdin, err = m.cmd.StdinPipe()
	if err != nil {
		return err
	}
	m.chWrite = make(chan []uint8, 1)
	m.stop = false

	go func() {
		out, _ := m.cmd.CombinedOutput()
		fmt.Printf("FFMPEG output:\n%v\n", string(out))
	}()

	go m.handleFrames(width, height)

	return nil
}

func (m FFMpeg) handleFrames(width, height int) {
	buf := new(bytes.Buffer)
	for {
		select {
		case pix := <-m.chWrite:
			stride := len(pix) / height
			rowLen := 4 * width
			for i := 0; i < len(pix); i += stride {
				if _, err := m.stdin.Write(pix[i : i+rowLen]); err != nil {
					break
				}
			}
			buf.Reset()
		}
	}
}

func (m *FFMpeg) Stop() {
	fmt.Print("chan", len(m.chWrite))
	m.stop = true
	m.stdin.Close()
}
