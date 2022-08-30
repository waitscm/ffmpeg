package ffmpeg

import "testing"

func TestTakeScreenShot(t *testing.T) {

	tests := []struct {
		name    string
		path    string
		out     string
		wantErr bool
	}{
		{
			"big buck",
			"./test/big_buck_bunny.mp4",
			"./test/big_buck_bunny.jpg",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := TakeScreenShot(tt.path, tt.out, 30); (err != nil) != tt.wantErr {
				t.Errorf("TakeScreenShot() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
