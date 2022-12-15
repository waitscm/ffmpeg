package ffmpeg

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
)

type (
	Stats struct {
		Streams Streams `json:"streams"`
		Format  Format  `json:"format"`
	}

	Streams []Stream

	Stream struct {
		CodecName   string `json:"codec_name"`
		CodecType   string `json:"codec_type"`
		Width       int    `json:"width"`
		Height      int    `json:"height"`
		DurationStr string `json:"duration"`
		BitRateStr  string `json:"bit_rate"`
	}

	Format struct {
		Name        string `json:"format_name"`
		DurationStr string `json:"duration"`
		BitRateStr  string `json:"bit_rate"`
		Tags        Tags   `json:"tags"`
	}

	Tags struct {
		Title string `json:"title"`
	}
)

const StreamTypeVideo = "video"

func GetVideoStats(path string) (Stats, error) {

	ospath := ensureVolPath(path)

	//command := "ffprobe -i 'ospath' -show_entries format=format_name,duration,bit_rate:format_tags=title:stream=codec_name,codec_type,width,height,duration,bit_rate -v quiet -of json"

	parts := []string{
		"ffprobe",
		"-i",
		ospath,
		"-show_entries",
		"format=format_name,duration,bit_rate:format_tags=title:stream=codec_name,codec_type,width,height,duration,bit_rate",
		"-v",
		"quiet",
		"-of",
		"json",
	}

	data, err := exec.Command(parts[0], parts[1:]...).Output()
	if err != nil {
		return Stats{}, fmt.Errorf("ffmpeg:GetVideoStats failed %v - %v: %w", ospath, string(data), err)
	}

	var s Stats
	err = json.Unmarshal(data, &s)
	if err != nil {
		return Stats{}, fmt.Errorf("ffmpeg:GetVideoStats failed unmarshal %v - %v: %w", ospath, string(data), err)
	}

	return s, nil
}

func (s Stream) IsVideo() bool {
	return s.CodecType == StreamTypeVideo
}

func (ss Streams) GetVideoStreamStats() (width, height int, bitRate int64) {

	for _, stream := range ss {
		if stream.IsVideo() {
			bitRate, _ := strconv.ParseInt(stream.BitRateStr, 10, 64)
			return stream.Width, stream.Height, bitRate
		}
	}
	return 0, 0, 0
}

func (st Stats) GetVideoStats() (width, height int, bitRate int64, duration float64) {
	width, height, bitRate = st.Streams.GetVideoStreamStats()
	duration, _ = strconv.ParseFloat(st.Format.DurationStr, 64)
	return
}

// GetTitle returns the tag title if present
func (st Stats) GetTitle() string {
	return st.Format.Tags.Title
}
