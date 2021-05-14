package measure

import (
	"reflect"
	"testing"
)

func TestSec_Tommssは再生時間をいいかんじのフォーマットにしてくれる(t *testing.T) {
	tests := []struct {
		name string
		s    Sec
		want string
	}{
		{"00:0s", Sec(0), "00:00"},
		{"00:ss", Sec(1), "00:01"},
		{"00:ss", Sec(59), "00:59"},
		{"0m:00", Sec(60), "01:00"},
		{"mm:00", Sec(720), "12:00"},
		{"0h:00:00", Sec(7200), "02:00:00"},
		{"hh:mm:ss", Sec(72072), "20:01:12"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Tommss(); got != tt.want {
				t.Errorf("Tommss() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadVideosは特定ディレクトリ直下のmp4をVideosに変換できる(t *testing.T) {
	type args struct {
		rootPath string
	}
	tests := []struct {
		name    string
		args    args
		want    Videos
		wantErr bool
	}{
		{"ex01",
			args{rootPath: "fixture"},
			Videos{
				Video{Path: "fixture/sample1.mp4", Sec: Sec(3)},
				Video{Path: "fixture/sample2.mp4", Sec: Sec(4)},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadVideos(tt.args.rootPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadVideos() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadVideos() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVideos_TotalSecは総再生時間を取得できる(t *testing.T) {
	tests := []struct {
		name string
		vs   Videos
		want Sec
	}{
		{"ex01",
			Videos{
				Video{Path: "fixture/sample1.mp4", Sec: Sec(3)},
				Video{Path: "fixture/sample2.mp4", Sec: Sec(4)},
			},
			Sec(7),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.vs.TotalSec(); got != tt.want {
				t.Errorf("TotalSec() = %v, want %v", got, tt.want)
			}
		})
	}
}
