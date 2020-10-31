package ffmpeg

import (
	"log"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
)

// GetVideoLength return the video length in seconds.
func GetVideoLength(path string) (float64, error) {

	ospath := filepath.FromSlash(path)
	command := "ffprobe -i '" + ospath + "' -show_entries format=duration -v quiet -of csv=p=0"

	lastQuote := rune(0)
	f := func(c rune) bool {
		switch {
		case c == lastQuote:
			lastQuote = rune(0)
			return false
		case lastQuote != rune(0):
			return false
		case unicode.In(c, unicode.Quotation_Mark):
			lastQuote = c
			return false
		default:
			return unicode.IsSpace(c)
		}
	}

	var parts []string
	preParts := strings.FieldsFunc(command, f)
	for i := range preParts {
		part := preParts[i]
		parts = append(parts, strings.Replace(part, "'", "", -1))
	}
	//parts = ["ffprobe", "-i", "'/media/Name of File.mp3'", "-show_entries", "format=duration", "-v", "quiet", "-of", "csv=p=0"]

	data, err := exec.Command(parts[0], parts[1:]...).Output()
	if err != nil {
		log.Println("ffmpeg:GetVideoLength failed ", ospath, err)
		return 0.0, err
	}

	stringVal := string(data)
	//fmt.Println("ffmpeg:GetVideoLength", ospath, stringVal)

	// strip \n \r characters
	stringVal = strings.TrimSuffix(stringVal, "\n")
	stringVal = strings.TrimSuffix(stringVal, "\r")

	length, err := strconv.ParseFloat(stringVal, 64)
	if err != nil {
		log.Println("ffmpeg:GetVideoLength can't parse float ", ospath, stringVal, err)
		return 0.0, err
	}
	return length, nil
}
