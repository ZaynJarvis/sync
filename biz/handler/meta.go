package handler

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"net/http"
	"sync/biz/dal/store"
	"sync/biz/dal/vod"
)

// VID gets meta info vid->urls
func VID(ctx context.Context, c *app.RequestContext) {
	vid := c.Param("vid")
	info, err := vod.PlayInfo(vid)
	if err != nil {
		c.JSON(http.StatusOK, utils.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, utils.H{"info": info})
}

func Videos(ctx context.Context, c *app.RequestContext) {
	list := store.List()
	for _, v := range list {
		vod.QueryTaskByID(v.ID) // update status
	}
	c.JSON(http.StatusOK, list)
}
