package handler

import (
	"context"
	"encoding/json"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/google/uuid"
	"net/http"
	"sync/biz/dal/store"
	"sync/biz/dal/vod"
)

// Sync a video from source_url to volc-vod
// 1. add a record (and return)
// 2. download video from source_url
// 3. upload to volc-vod
// 4. update sync process
func Sync(ctx context.Context, c *app.RequestContext) {
	body, _ := c.Body()
	var urls []string
	json.Unmarshal(body, &urls)

	var videos []store.VideoInfo
	jobInfo := vod.UploadByUrl(urls)
	for _, job := range jobInfo {
		info := store.VideoInfo{
			ID:        uuid.NewString(),
			Status:    store.Pending,
			SourceURL: job.SourceURL,
			JobID:     job.ID,
		}
		store.Put(info)
		videos = append(videos, info)
	}

	c.JSON(http.StatusOK, videos)
}

func SyncStatus(ctx context.Context, c *app.RequestContext) {
	id := c.Param("id")
	status := vod.QueryTaskByID(id)
	c.JSON(http.StatusOK, utils.H{"status": status})
}
