package ffmpeg

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func GetVideoLength(path string) (float64, error) {

	ospath := ensureVolPath(path)

	//command := "ffprobe -i " + ospath + " -show_entries format=duration -v quiet -of csv=p=0"

	parts := []string{
		"ffprobe",
		"-i",
		ospath, //strings.Replace(ospath, "'", "\\'", -1),
		"-show_entries",
		"format=duration",
		"-v",
		"quiet",
		"-of",
		"csv=p=0",
	}

	data, err := exec.Command(parts[0], parts[1:]...).Output()
	if err != nil {
		return 0.0, fmt.Errorf("ffmpeg:GetVideoLength failed %v - %v: %w", ospath, string(data), err)
	}

	stringVal := string(data)

	// strip \n \r characters
	stringVal = strings.TrimSuffix(stringVal, "\n")
	stringVal = strings.TrimSuffix(stringVal, "\r")

	length, err := strconv.ParseFloat(stringVal, 64)
	if err != nil {
		return 0.0, fmt.Errorf("ffmpeg:GetVideoLength can't parse float %v - %v: %w", ospath, stringVal, err)
	}
	return length, nil
}

// if there's an apostrophe ToSlash doesn't always use the C:/ format which confuses ffprobe
//  in windows.
func ensureVolPath(path string) string {
	ospath := filepath.ToSlash(path)
	if runtime.GOOS == "windows" && len(ospath) > 1 && ospath[0] == '/' {
		vol := fmt.Sprintf("%v:", string(ospath[1]))
		ospath = vol + ospath[2:]
	}
	return ospath
}
