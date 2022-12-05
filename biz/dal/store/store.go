// Package store keeps track of all video syncing
// more sophisticated app need to use a database
package store

import "sync"

type Status int8

const (
	Unknown Status = iota
	Pending
	Success
	Failed
)

type VideoInfo struct {
	VID       string
	Status    Status
	SourceURL string
}

var (
	mu    sync.Mutex
	cache = make(map[string]VideoInfo)
)

func Get(vid string) VideoInfo {
	mu.Lock()
	defer mu.Unlock()

	return cache[vid]
}

func Put(info VideoInfo) {
	mu.Lock()
	defer mu.Unlock()

	cache[info.VID] = info
}
