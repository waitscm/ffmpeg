package ffmpeg

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type (
	Stats struct {
		Streams Streams `json:"streams"`
		Format  Format  `json:"format"`
	}

	Streams []Stream

	Stream struct {
		CodecName           string `json:"codec_name"`
		CodecType           string `json:"codec_type"`
		PixelFormat         string `json:"pix_fmt"`
		AvgFrameRateRaw     string `json:"avg_frame_rate"`
		AvgFrameRate        float64
		BitsPerRawSampleStr string `json:"bits_per_raw_sample"`
		BitsPerRawSample    int
		Width               int    `json:"width"`
		Height              int    `json:"height"`
		DurationStr         string `json:"duration"`
		Duration            float64
		BitRateStr          string `json:"bit_rate"`
		BitRate             int64
	}

	Format struct {
		Name        string `json:"format_name"`
		DurationStr string `json:"duration"`
		Duration    float64
		BitRateStr  string `json:"bit_rate"`
		BitRate     int64
		Tags        Tags `json:"tags"`
	}

	Tags struct {
		Title string `json:"title"`
	}
)

const StreamTypeVideo = "video"

func GetVideoStats(path string) (Stats, error) {

	ospath := ensureVolPath(path)

	//command := "ffprobe -i 'ospath' -show_entries format=format_name,duration,bit_rate:format_tags=title:stream=avg_frame_rate,pix_fmt,bits_per_raw_sample,codec_name,codec_type,width,height,duration,bit_rate -v quiet -of json"

	parts := []string{
		"ffprobe",
		"-i",
		ospath,
		"-show_entries",
		"format=format_name,duration,bit_rate:format_tags=title:stream=avg_frame_rate,pix_fmt,bits_per_raw_sample,codec_name,codec_type,width,height,duration,bit_rate",
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

	s.calc()

	return s, nil
}

func (s Stream) IsVideo() bool {
	return s.CodecType == StreamTypeVideo
}

func (st Stats) GetVideoStream() (Stream, bool) {

	for _, stream := range st.Streams {
		if stream.IsVideo() {
			return stream, true
		}
	}
	return Stream{}, false
}

// GetTitle returns the tag title if present
func (st Stats) GetTitle() string {
	return st.Format.Tags.Title
}

func (f *Format) calc() {
	f.BitRate, _ = strconv.ParseInt(f.BitRateStr, 10, 64)
	f.Duration, _ = strconv.ParseFloat(f.DurationStr, 64)
}

func (s *Stream) calc() {
	s.BitRate, _ = strconv.ParseInt(s.BitRateStr, 10, 64)
	s.Duration, _ = strconv.ParseFloat(s.DurationStr, 64)
	s.BitsPerRawSample, _ = strconv.Atoi(s.BitsPerRawSampleStr)
	if s.AvgFrameRateRaw != "" {
		fr := strings.Split(s.AvgFrameRateRaw, "/")
		if len(fr) == 2 {
			n, _ := strconv.Atoi(fr[0])
			d, _ := strconv.Atoi(fr[1])
			if d != 0 {
				s.AvgFrameRate = float64(n) / float64(d)
			}
		}
	}
}

func (st *Stats) calc() {

	st.Format.calc()
	for i := range st.Streams {
		st.Streams[i].calc()
	}
}
