// Package vod handles volc-vod upload and play_info
package vod

import (
	"fmt"
	"github.com/volcengine/volc-sdk-golang/base"
	"github.com/volcengine/volc-sdk-golang/service/vod"
	"github.com/volcengine/volc-sdk-golang/service/vod/models/business"
	"github.com/volcengine/volc-sdk-golang/service/vod/models/request"
	"log"
	"os"
	"sync/biz/dal/store"
)

type JobInfo struct {
	ID        string
	SourceURL string
}

func PlayInfo(vid string) ([]string, error) {
	instance := vod.NewInstance()
	// get your own AK & SK from volc-vod console
	instance.SetCredential(base.Credentials{
		AccessKeyID:     os.Getenv("VOD_AK"),
		SecretAccessKey: os.Getenv("VOD_SK"),
	})
	fmt.Println(vid)
	resp, _, err := instance.GetPlayInfo(&request.VodGetPlayInfoRequest{Vid: vid})
	if err != nil {
		return nil, err
	}
	fmt.Println(resp.String())

	var urls []string
	for _, v := range resp.Result.PlayInfoList {
		urls = append(urls, v.MainPlayUrl)
	}
	return urls, nil
}

func UploadByUrl(urls []string) []JobInfo {
	instance := vod.NewInstance()
	// get your own AK & SK from volc-vod console
	instance.SetCredential(base.Credentials{
		AccessKeyID:     os.Getenv("VOD_AK"),
		SecretAccessKey: os.Getenv("VOD_SK"),
	})

	var bizURLSet []*business.VodUrlUploadURLSet
	for _, url := range urls {
		bizURLSet = append(bizURLSet, &business.VodUrlUploadURLSet{SourceUrl: url})
	}

	req := &request.VodUrlUploadRequest{
		SpaceName: os.Getenv("VOD_SPACE"),
		URLSets:   bizURLSet,
	}

	resp, _, err := instance.UploadMediaByUrl(req)
	fmt.Println(resp.String())
	if err != nil {
		return nil
	}
	if resp.Result == nil {
		return nil
	}
	var jobInfo []JobInfo
	for _, kv := range resp.Result.Data {
		jobInfo = append(jobInfo, JobInfo{
			SourceURL: kv.SourceUrl,
			ID:        kv.JobId,
		})
	}
	return jobInfo
}

const (
	Unknown = iota
	Pending
	Success
	Failed
	Published
)

func QueryTaskByID(id string) store.Status {
	videoInfo := store.Get(id)
	if videoInfo.Status == store.Pending {
		var status store.Status
		vodStatus, vid := queryTaskByJobID(videoInfo.JobID)
		switch vodStatus {
		case Pending:
			status = store.Pending
		case Success:
			videoInfo.VID = vid
			status = store.Success
		case Failed:
			status = store.Failed
		case Unknown:
			status = store.Unknown
		}
		videoInfo.Status = status
		store.Put(videoInfo)
	}
	return videoInfo.Status
}

func queryTaskByJobID(jobID string) (int, string) {
	instance := vod.NewInstance()
	instance.SetCredential(base.Credentials{
		AccessKeyID:     os.Getenv("VOD_AK"),
		SecretAccessKey: os.Getenv("VOD_SK"),
	})

	req := &request.VodQueryUploadTaskInfoRequest{
		JobIds: jobID,
	}

	resp, _, err := instance.QueryUploadTaskInfo(req)
	fmt.Println(resp.String())
	if err != nil {
		return Unknown, ""
	}
	if resp.Result == nil {
		return Unknown, ""
	}
	if len(resp.Result.Data.MediaInfoList) == 0 {
		return Unknown, ""
	}
	if resp.Result.Data.MediaInfoList[0].JobId != jobID {
		log.Println("jobID changed")
		return Unknown, ""
	}
	switch resp.Result.Data.MediaInfoList[0].State {
	case "success":
		return Success, resp.Result.Data.MediaInfoList[0].Vid
	case "failed":
		return Failed, ""
	case "processing":
		return Pending, ""
	default:
		return Unknown, ""
	}
}
