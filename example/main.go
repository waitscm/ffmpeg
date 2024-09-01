package main

import (
	"fmt"
	"os"
	"time"

	"github.com/waitscm/ffmpeg"
)

func intPtr(i int) *int {
	return &i
}

func floatPtr(i float32) *float32 {
	return &i
}

func stringPtr(i string) *string {
	return &i
}

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

	dur = dur * 1000 / 4.0

	start = time.Now()
	err = ffmpeg.TakeScreenShotMS(filename, "screenshot.jpg", int64(dur))
	if err != nil {
		fmt.Println("TakeScreenShot", err)
		os.Exit(1)
	}
	fmt.Println("TakeScreenShot elapsed", time.Since(start))

	tests := []struct {
		name string
		f    ffmpeg.ScreenFilter
	}{
		{
			name: "320x240",
			f: ffmpeg.ScreenFilter{
				Width:  intPtr(320),
				Height: intPtr(240),
			},
		},
		{
			name: "320x-1",
			f: ffmpeg.ScreenFilter{
				Width: intPtr(320),
			},
		},
		{
			name: "samesize",
			f:    ffmpeg.ScreenFilter{},
		},
		{
			name: "half",
			f: ffmpeg.ScreenFilter{
				ReductionFactor: floatPtr(2),
			},
		},
		{
			name: "double",
			f: ffmpeg.ScreenFilter{
				ReductionFactor: floatPtr(.5),
			},
		},
		{
			name: "200x200pad",
			f: ffmpeg.ScreenFilter{
				Width:  intPtr(200),
				Height: intPtr(200),
				Pad:    true,
			},
		},
		{
			name: "320x240hflip",
			f: ffmpeg.ScreenFilter{
				Width:          intPtr(320),
				Height:         intPtr(240),
				HorizontalFlip: true,
			},
		},
		{
			name: "320x-1hflip",
			f: ffmpeg.ScreenFilter{
				Width:          intPtr(320),
				HorizontalFlip: true,
			},
		},
		{
			name: "halfhflip",
			f: ffmpeg.ScreenFilter{
				ReductionFactor: floatPtr(2),
				HorizontalFlip:  true,
			},
		},
		{
			name: "doublehflip",
			f: ffmpeg.ScreenFilter{
				ReductionFactor: floatPtr(.5),
				HorizontalFlip:  true,
			},
		},
		{
			name: "200x200padhflip",
			f: ffmpeg.ScreenFilter{
				Width:          intPtr(200),
				Height:         intPtr(200),
				Pad:            true,
				HorizontalFlip: true,
			},
		},
		{
			name: "rawhflip",
			f: ffmpeg.ScreenFilter{
				RawFilter: stringPtr("hflip"),
			},
		},
	}

	for _, t := range tests {
		start = time.Now()
		out := fmt.Sprintf("%v.jpg", t.name)

		ffmpeg.TakeScreenShotFilter(filename, out, int64(dur), t.f)

		fmt.Println("TakeScreenShotFilter elapsed", t.name, time.Since(start))
	}

	for i := 1; i <= 32; i += 2 {
		start = time.Now()
		out := fmt.Sprintf("quality_%v.jpg", i)

		ffmpeg.TakeScreenShotFilter(filename, out, int64(dur), ffmpeg.ScreenFilter{Quality: intPtr(i)})

		fmt.Println("TakeScreenShotFilter elapsed quality", i, time.Since(start))
	}

	start = time.Now()
	stats, err := ffmpeg.GetVideoStats(filename)
	if err != nil {
		fmt.Println("GetVideoStats", err)
	}
	fmt.Println("GetVideoStats elapsed", time.Since(start))
	vs, has := stats.GetVideoStream()
	fmt.Println("VideoStats", "\n\tfilename", filename, "\n\ttitle:", stats.Format.Tags.Title, "\n\thas video stream:", has, "\n\tw", vs.Width, "x h", vs.Height,
		"\n\tbitrate:", vs.BitRate, "\n\tduration s:", vs.Duration, "\n\tpixel format:", vs.PixelFormat, "\n\tcodec:", vs.CodecName,
		"\n\tbits per raw sample:", vs.BitsPerRawSample, "\n\tavg frame rate:", vs.AvgFrameRate, "raw:", vs.AvgFrameRateRaw)
}
