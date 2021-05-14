package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"work-measure-mp4-duration/measure"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("引数はディレクトリ1つな")
		os.Exit(1)
	}

	targetPath := os.Args[1]

	videos, err := measure.ReadVideos(targetPath)
	if err != nil {
		log.Fatal(err)
	}

	if videos == nil {
		log.Fatalf("%s直下にはmp4ファイルは1つもなかったよ\n", targetPath)
	}

	for _, video := range videos {
		relPath, err := filepath.Rel(targetPath, video.Path)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\t%s\n", video.Sec.Tommss(), relPath)
	}

	fmt.Printf("合計時間 %s\n", videos.TotalSec().Tommss())
}
