// Package vod handles volc-vod upload and play_info
package vod

import (
	"fmt"
	"log"
	"math/rand"
)

func InitUpload() string {
	return fmt.Sprintf("vid:%d", rand.Int())
}

func Upload(vid string, data []byte) error {
	log.Println(vid, len(data))
	return nil
}

func PlayInfo(vid string) []string {
	log.Println(vid)
	return nil
}
