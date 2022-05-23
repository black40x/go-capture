package pkg

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
)

type FFMpeg struct {
}

func readStuff(scanner *bufio.Scanner) {
	for scanner.Scan() {
		fmt.Println("Performed Scan")
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}

func (m FFMpeg) Version() (string, error) {
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
