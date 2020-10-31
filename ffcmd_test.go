package ffmpeg

import (
	"log"
	"testing"
)

func TestGetVideoLength(t *testing.T) {
	mp4 := "test/big_buck_bunny.mp4"
	avi := "test/drop.avi"
	flv := "test/small.flv"
	notvid := "test/notvideo.mp4"
	noexist := "test/tttt.mp4"

	len, err := GetVideoLength(mp4)
	if err != nil {
		t.Error("failed mp4", err)
	}
	log.Println("mp4 length", len)

	len, err = GetVideoLength(avi)
	if err != nil {
		t.Error("failed avi", err)
	}
	log.Println("avi length", len)

	len, err = GetVideoLength(flv)
	if err != nil {
		t.Error("failed flv", err)
	}
	log.Println("flv length", len)

	_, err = GetVideoLength(notvid)
	if err == nil {
		t.Error("failed to err on file that's not a video")
	}

	_, err = GetVideoLength(noexist)
	if err == nil {
		t.Error("failed to err on file that doesn't exist")
	}
}
