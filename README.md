# ffmpeg

Requires ffmpeg installed.

## Video length

Get the video length of a file.

`length, err := GetVideoLength("/path/to/video.mp4")`

## Take a screenshot

Pass in the video and output screenshot with second mark to take screen at.

`err := TakeScreenShot("/path/to/video.mp4", "/path/to/out.jpg", 25)`