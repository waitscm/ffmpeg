package main

import (
	"fmt"
	"os"
	"time"

	"github.com/waitscm/ffmpeg"
)

func main() {
	if len(os.Args) == 0 {
		fmt.Println("./example <video file>")
		os.Exit(0)
	}
	filename := os.Args[1]
	start := time.Now()
	dur, err := ffmpeg.GetVideoLength(filename)
	if err != nil {
		fmt.Println("error getting length", err)
		os.Exit(1)
	}
	fmt.Println("GetVideoLength elapsed", time.Since(start), filename, "duration", dur)

	if dur > 20.0 {
		dur -= 7.0
	} else {
		dur = dur / 2.0
	}

	start = time.Now()
	err = ffmpeg.TakeScreenShot(filename, "screenshot.jpg", int(dur))
	if err != nil {
		fmt.Println("TakeScreenShot", err)
	}
	fmt.Println("TakeScreenShot elapsed", time.Since(start))

	start = time.Now()
	stats, err := ffmpeg.GetVideoStats(filename)
	if err != nil {
		fmt.Println("GetVideoStats", err)
	}
	fmt.Println("GetVideoStats elapsed", time.Since(start), stats)
	w, h, br, dur := stats.GetVideoStats()
	fmt.Println("width:", w, "height:", h, "bitrate", br, "duration", dur)
}
