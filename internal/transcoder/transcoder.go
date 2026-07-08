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
	Input  string `json:"input" validate:"required"`
	Output string `json:"output" validate:"required"`
}

type Progress struct {
	Frame   string
	FPS     string
	Bitrate string
	Outtime string
	Speed   string
}

func updateProgress(progress *Progress, key string, val string) {

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

		if len(split) < 2 {
			continue
		}

		key := split[0]
		val := split[1]

		updateProgress(&progress, key, val)

		if key == "progress" {
			select {
			case progresschan <- progress:
				progress = Progress{}
			case <-ctx.Done():
				return ctx.Err()
			}
		}

	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner failed, %w", err)
	}

	if err := cmd.Wait(); err != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return fmt.Errorf("ffmpeg process failed, %w", err)
	}

	return nil
}
