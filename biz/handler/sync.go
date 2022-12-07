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

const SyncBatch = 20

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
	for _, batch := range batched(urls, SyncBatch) {
		jobInfo := vod.UploadByUrl(batch)
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
	}

	c.JSON(http.StatusOK, videos)
}

func SyncStatus(ctx context.Context, c *app.RequestContext) {
	id := c.Param("id")
	status := vod.QueryTaskByID(id)
	c.JSON(http.StatusOK, utils.H{"status": status})
}

func batched(urls []string, batchLimit int) (batchedURL [][]string) {
	for i := 0; i <= len(urls)/batchLimit; i++ {
		end := (i + 1) * batchLimit
		if len(urls) < end {
			end = len(urls)
		}
		batchedURL = append(batchedURL, urls[i*batchLimit:end])
	}
	return
}
