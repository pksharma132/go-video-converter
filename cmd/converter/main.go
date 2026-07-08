package main

import (
	"context"
	"fmt"
	"go-video-converter/internal/transcoder"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func run() error {
	args := os.Args

	if len(args) != 3 {
		return fmt.Errorf("usage converter <input> <output>")
	}

	input := args[1]
	output := args[2]

	t := transcoder.New()
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer cancel()
	progressChan := make(chan transcoder.Progress)
	errorChan := make(chan error, 1)

	go func() {
		errorChan <- t.Convert(ctx, transcoder.ConvertParam{Input: input, Output: output}, progressChan)
	}()

	for {
		select {
		case p, ok := <-progressChan:
			if !ok {
				progressChan = nil
				continue

			}

			fmt.Println(p)

		case err := <-errorChan:
			if err != nil {
				return fmt.Errorf("fatal error, %w", err)
			}

			return nil

		}
	}
}

func main() {

	if err := run(); err != nil {
		log.Fatal(err)
	}
}
