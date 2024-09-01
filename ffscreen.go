package ffmpeg

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
)

type ScreenFilter struct {
	// Width if set will set output width. Otherwise if Height set and not Width will use -1.
	Width *int
	// Height if set will set output height. Otherwise if Width set and not Height will use -1.
	Height *int
	// ReductionFactor if set will reduce height and width by this factor
	ReductionFactor *float32
	// Pad if set will force original vid aspect ratio to be maintained and pad.
	Pad bool
	// RawFilter will be passed directly to the -vf argument in ffmpeg and ignore all other settings
	RawFilter *string
	// HorizontalFlip will flip the image over Y-axis
	HorizontalFlip bool
	// Quality sets the output quality 0 being lossless and 51 being worst quality. defaults to 32
	Quality *int
}

func TakeScreenShot(inputPath string, outputPath string, seekSeconds int) error {
	//ffmpeg.exe -ss 00:00:30 -i ./test/big_buck_bunny.mp4 -vframes 1 -q:v 31 output.jpg
	ospath := ensureVolPath(inputPath)
	outPath := ensureVolPath(outputPath)

	seek := fmt.Sprintf("%02d:%02d:%02d", seekSeconds/3600, (seekSeconds%3600)/60, seekSeconds%60)

	return takeSShot(ospath, outPath, seek, nil)
}

func TakeScreenShotMS(inputPath string, outputPath string, seekMilliseconds int64) error {
	//ffmpeg.exe -ss 5023ms -i ./test/big_buck_bunny.mp4 -vframes 1 -q:v 31 output.jpg
	ospath := ensureVolPath(inputPath)
	outPath := ensureVolPath(outputPath)

	seek := fmt.Sprintf("%dms", seekMilliseconds)

	return takeSShot(ospath, outPath, seek, nil)
}

func TakeScreenShotFilter(inputPath string, outputPath string, seekMilliseconds int64, f ScreenFilter) error {
	//ffmpeg.exe -ss 5023ms -i ./test/big_buck_bunny.mp4 -vframes 1 -q:v 31 output.jpg
	ospath := ensureVolPath(inputPath)
	outPath := ensureVolPath(outputPath)

	seek := fmt.Sprintf("%dms", seekMilliseconds)

	return takeSShot(ospath, outPath, seek, &f)
}

func (s *ScreenFilter) String() string {
	if s == nil {
		return ""
	}

	if s.RawFilter != nil {
		return *s.RawFilter
	}

	if s.ReductionFactor != nil {
		if s.HorizontalFlip {
			return fmt.Sprintf("scale=iw/%v:ih/%v,hflip", *s.ReductionFactor, *s.ReductionFactor)
		}
		return fmt.Sprintf("scale=iw/%v:ih/%v", *s.ReductionFactor, *s.ReductionFactor)
	}

	height := -1
	width := -1

	if s.Height != nil {
		height = *s.Height
	}
	if s.Width != nil {
		width = *s.Width
	}
	if s.Pad && width > -1 && height > -1 {
		if s.HorizontalFlip {
			return fmt.Sprintf("scale=%v:%v:force_original_aspect_ratio=decrease,pad=%v:%v:(ow-iw)/2:(oh-ih)/2,hflip", width, height, width, height)
		}
		return fmt.Sprintf("scale=%v:%v:force_original_aspect_ratio=decrease,pad=%v:%v:(ow-iw)/2:(oh-ih)/2", width, height, width, height)
	}

	if s.HorizontalFlip {
		return fmt.Sprintf("scale=%v:%v,hflip", width, height)
	}
	return fmt.Sprintf("scale=%v:%v", width, height)
}

func takeSShot(in, out, seek string, f *ScreenFilter) error {
	cmdStr := []string{"-hide_banner", "-loglevel", "error", "-y", "-ss", seek, "-i", in, "-frames:v", "1", "-q:v"}
	quality := 32
	if f != nil {
		if f.Quality != nil {
			quality = *f.Quality
		}
	}
	cmdStr = append(cmdStr, fmt.Sprintf("%v", quality))

	if f != nil {
		cmdStr = append(cmdStr, "-vf", f.String())
	}
	cmdStr = append(cmdStr, out)
	cmd := exec.Command("ffmpeg", cmdStr...)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	slurp, _ := ioutil.ReadAll(stderr)

	if len(slurp) > 0 {
		return errors.New(string(slurp))
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}
