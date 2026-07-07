package main

import (
	"context"
	"go-video-converter/internal/transcoder"
	"log"
)

func main() {
	t := transcoder.New()
	ctx := context.Background()
	progresschan := make(chan transcoder.Progress)
	// errorchan := make(chan error, 1)

	// go func() {
	err := t.Convert(ctx, transcoder.ConvertParam{Input: "lotm.webm", Output: "output.webm"}, progresschan)
	if err != nil {
		log.Fatal(err)
	}
	// errorchan <- err
	// }()

}
