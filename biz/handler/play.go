package handler

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"net/http"
	"sync/biz/dal/store"
	"sync/biz/dal/vod"
)

// Play a video based on vid, if video is under sync, use source url, otherwise use volc-vod
// 1. get record
// 2. check sync status
// 3.1 if sync pending, return source_url
// 3.2 if sync success, depend on play_param return volc-vod transcode_url
func Play(ctx context.Context, c *app.RequestContext) {
	vid := c.Param("vid")
	videoInfo := store.Get(vid)

	switch videoInfo.Status {
	case store.Unknown, store.Failed:
		c.JSON(http.StatusNotFound, utils.H{"message": "vid not found"})
	case store.Pending:
		c.JSON(http.StatusOK, utils.H{"url": videoInfo.SourceURL})
	case store.Success:
		urls := vod.PlayInfo(vid)
		if len(urls) == 0 {
			c.JSON(http.StatusNotFound, utils.H{"message": "vid not found"})
			return
		}
		c.JSON(http.StatusOK, utils.H{"url": TranscodeDecider( /* transcodeParam, */ urls)})
	}
}

func TranscodeDecider( /* transcodeParam, */ urls []string) string {
	return urls[0]
}