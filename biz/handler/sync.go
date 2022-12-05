package handler

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"log"
	"net/http"
	"sync/biz/dal/source"
	"sync/biz/dal/store"
	"sync/biz/dal/vod"
)

// Sync a video from source_url to volc-vod
// 1. add a record (and return)
// 2. download video from source_url
// 3. upload to volc-vod
// 4. update sync process
func Sync(ctx context.Context, c *app.RequestContext) {
	sourceURL := c.Query("source_url")
	if sourceURL == "" {
		c.JSON(http.StatusBadRequest, utils.H{"message": "need source_url in query"})
		return
	}

	vid := vod.InitUpload()
	store.Put(store.VideoInfo{
		VID:       vid,
		Status:    store.Pending,
		SourceURL: sourceURL,
	})
	c.JSON(http.StatusOK, utils.H{"message": "sync started"})

	// the sync process can be in background, a cronjob, a go routine
	// can block for response, wait for callback, and update when use
	// here shows the most simple way (not necessary robust)
	data, err := source.Get(sourceURL)
	if err != nil {
		store.Put(store.VideoInfo{VID: vid, Status: store.Failed})
		log.Printf("sourceURL %s not found\n", sourceURL)
		return
	}
	log.Printf("sourceURL %s downloaded, data len: %d\n", sourceURL, len(data))

	err = vod.Upload(vid, data)
	if err != nil {
		store.Put(store.VideoInfo{VID: vid, Status: store.Failed}) // may use source_url to fallback
		log.Printf("vid %s upload failed, err: %v\n", vid, err)
		return
	}
	log.Printf("vid %s uploaded\n", vid)

	store.Put(store.VideoInfo{VID: vid, Status: store.Success})
}

func SyncStatus(ctx context.Context, c *app.RequestContext) {
	vid := c.Param("vid")

	videoInfo := store.Get(vid)
	c.JSON(http.StatusOK, utils.H{"status": videoInfo.Status})
}
