package measure

import (
	"errors"
	"fmt"
	"github.com/abema/go-mp4"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type Sec int

func (s Sec) Tommss() string {
	ss := s % 60
	mm := (s / 60) % 60
	hh := s / 60 / 60

	if 0 <= s && s < 3600 {
		return fmt.Sprintf("%02d:%02d", mm, ss)
	} else if 3600 <= s {
		return fmt.Sprintf("%02d:%02d:%02d", hh, mm, ss)
	} else {
		return ""
	}
}

type Video struct {
	Path string
	Sec  Sec
}

type Videos []Video

func (vs Videos) TotalSec() Sec {
	var total Sec

	for _, v := range vs {
		total += v.Sec
	}

	return total
}

// rootPathの直下にある全てのmp4ファイルの情報をファイル名昇順でまとめる
func ReadVideos(rootPath string) (Videos, error) {
	var videos Videos

	f, err := os.Open(rootPath)
	if err != nil {
		return nil, err
	}

	filenames, err := f.Readdirnames(-1)
	if err != nil {
		return nil, err
	}

	sort.Strings(filenames)

	for _, filename := range filenames {
		info, _ := os.Stat(filepath.Join(rootPath, filename))
		if info.IsDir() {
			continue
		}

		file, err := os.Open(filepath.Join(rootPath, filename))
		if err != nil {
			return nil, err
		}

		if filepath.Ext(file.Name()) == ".mp4" {
			_, err = mp4.ReadBoxStructure(file, func(h *mp4.ReadHandle) (interface{}, error) {
				if h.BoxInfo.Type == mp4.BoxTypeMvhd() {
					sec, err := parseSec(h)
					if err != nil {
						return nil, err
					}
					video := Video{Path: file.Name(), Sec: sec}

					videos = append(videos, video)
				}

				if h.BoxInfo.IsSupportedType() {
					return h.Expand()
				}

				return nil, nil
			})

			if err != nil {
				return nil, err
			}
		}
	}

	return videos, err
}

// MVHDボックスから再生時間(秒)を取得する
// こんな感じの文字列を強引にParse
// Version=0 Flags=0x000000 CreationTimeV0=3703726038 ModificationTimeV0=3703726038 Timescale=30000 DurationV0=2473200 Rate=1 Volume=256 Matrix=[0x10000, 0x0, 0x0, 0x0, 0x10000, 0x0, 0x0, 0x0, 0x40000000] PreDefined=[0, 0, 0, 0, 0, 0] NextTrackID=3
func parseSec(h *mp4.ReadHandle) (Sec, error) {
	box, _, err := h.ReadPayload()
	if err != nil {
		return 0, err
	}

	str, err := mp4.Stringify(box, mp4.Context{})
	if err != nil {
		return 0, err
	}

	if strings.Index(str, "Timescale") == -1 || strings.Index(str, "DurationV0") == -1 {
		return 0, errors.New("なんかおかしいよ")
	}

	split := strings.Split(str, " ")
	duration, err := strconv.Atoi(strings.Split(split[5], "=")[1])
	if err != nil {
		return 0, err
	}

	timescale, err := strconv.Atoi(strings.Split(split[4], "=")[1])
	if err != nil {
		return 0, err
	}

	return Sec(duration / timescale), nil
}
