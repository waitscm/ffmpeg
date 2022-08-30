package ffmpeg

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
)

func TakeScreenShot(inputPath string, outputPath string, seekSeconds int) error {
	//ffmpeg.exe -ss 00:00:30 -i ./test/big_buck_bunny.mp4 -vframes 1 -q:v 31 output.jpg
	ospath := ensureVolPath(inputPath)
	outPath := ensureVolPath(outputPath)

	seek := fmt.Sprintf("%02d:%02d:%02d", seekSeconds/3600, seekSeconds/60, seekSeconds%60)

	cmd := exec.Command("ffmpeg", "-hide_banner", "-loglevel", "error", "-y", "-ss", seek, "-i", ospath, "-frames:v", "1", "-q:v", "32", outPath)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	slurp, _ := ioutil.ReadAll(stderr)
	log.Printf("output%s\n", slurp)

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}

	return nil
}
