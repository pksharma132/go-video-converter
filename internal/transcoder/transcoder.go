package transcoder

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"strings"
)

type Transcoder struct {
}

func New() *Transcoder {
	return &Transcoder{}
}

type ConvertParam struct {
	Input  string
	Output string
}

type Progress struct {
	Frame   string
	FPS     string
	Bitrate string
	Outtime string
	Speed   string
}

func (t *Transcoder) Convert(ctx context.Context, param ConvertParam, progresschan chan<- Progress) error {

	defer close(progresschan)

	cmd := exec.CommandContext(ctx, "ffmpeg", "-progress", "pipe:1", "-nostats", "-i", param.Input, param.Output)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get the stdout pipe, %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start the command, %w", err)
	}

	scanner := bufio.NewScanner(stdout)

	progress := Progress{}

	for scanner.Scan() {
		t := scanner.Text()
		split := strings.SplitN(t, "=", 2)
		var key string
		var val string

		for i, a := range split {
			if i == 0 {
				key = a
			} else {
				val = a
			}
		}


		if key == "frame" {
			progress.Frame = val
		}

		if key == "bitrate" {
			progress.Bitrate = val
		}

		if key == "fps" {
			progress.FPS = val
		}

		if key == "out_time" {
			progress.Outtime = val
		}

		if key == "speed" {
			progress.Speed = val
		}

		if key == "progress" && val == "continue" {
			fmt.Println("progress", progress)
		}

	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner failed, %w", err)
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("ffmpeg process failed, %w", err)
	}

	return nil
}
