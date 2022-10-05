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

	if dur > 20.0 {
		dur -= 7.0
	} else {
		dur = dur / 2.0
	}

	err = ffmpeg.TakeScreenShot(filename, "screenshot.jpg", int(dur))
	if err != nil {
		fmt.Println("TakeScreenShot", err)
	}
}
