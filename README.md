# ffmpeg

Requires ffmpeg installed.

## Video Stats

Get information about the video

```
stats, _ := ffmpeg.GetVideoStats(filename)
vs, has := stats.GetVideoStream()
fmt.Println("VideoStats", "\n\tfilename", filename, "\n\ttitle:", stats.Format.Tags.Title, "\n\thas video stream:", has, "\n\tw", vs.Width, "x h", vs.Height,
		"\n\tbitrate:", vs.BitRate, "\n\tduration s:", vs.Duration, "\n\tpixel format:", vs.PixelFormat, "\n\tcodec:", vs.CodecName,
		"\n\tbits per raw sample:", vs.BitsPerRawSample, "\n\tavg frame rate:", vs.AvgFrameRate, "raw:", vs.AvgFrameRateRaw)
```

## Video length

Get the video length of a file.

`length, err := ffmpeg.GetVideoLength("/path/to/video.mp4")`

## Take a screenshot

Pass in the video and output screenshot with second mark to take screen at.

`err := TakeScreenShot("/path/to/video.mp4", "/path/to/out.jpg", 25)`