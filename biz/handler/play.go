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
	id := c.Param("id")
	status := vod.QueryTaskByID(id)
	videoInfo := store.Get(id)
	switch status {
	case store.Unknown, store.Failed: // depend on use case, can fallback to source url
		c.JSON(http.StatusNotFound, utils.H{"message": "vid upload failed"})
	case store.Pending:
		c.JSON(http.StatusOK, utils.H{"url": videoInfo.SourceURL, "type": "source_url", "reason": "uploading"})
	case store.Success:
		urls, err := vod.PlayInfo(videoInfo.VID)
		if err != nil {
			c.JSON(http.StatusNotFound, utils.H{"message": err.Error()})
			return
		}
		if len(urls) == 0 {
			c.JSON(http.StatusOK, utils.H{"url": videoInfo.SourceURL, "type": "source_url", "reason": "not published"})
			return
		}
		c.JSON(http.StatusOK, utils.H{"url": TranscodeDecider( /* transcodeParam, */ urls), "type": "volc_video"})
	}
}

func TranscodeDecider( /* transcodeParam, */ urls []string) string {
	return urls[0]
}
