package main

import (
	"fmt"
	"os"

	"github.com/waitscm/ffmpeg"
)

func main() {
	if len(os.Args) == 0 {
		fmt.Println("./example <video file>")
		os.Exit(0)
	}
	filename := os.Args[1]
	dur, err := ffmpeg.GetVideoLength(filename)
	if err != nil {
		fmt.Println("error getting length", err)
		os.Exit(1)
	}
	fmt.Println(filename, "duration", dur)
}
