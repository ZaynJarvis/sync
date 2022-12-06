// Package store keeps track of all video syncing
// more sophisticated app need to use a database
package store

import (
	"encoding/json"
	bolt "go.etcd.io/bbolt"
	"os"
	"time"
)

var (
	db, _  = bolt.Open(os.Getenv("VOD_SPACE")+".db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	bucket = []byte("VideoInfo")
)

func init() {
	db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists(bucket)
		return nil
	})
}

type Status int8

const (
	Unknown Status = iota
	Pending
	Success
	Failed
)

type VideoInfo struct {
	ID        string
	Status    Status
	SourceURL string
	JobID     string
	VID       string
}

func Get(id string) VideoInfo {
	var info VideoInfo
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		json.Unmarshal(b.Get([]byte(id)), &info)
		return nil
	})

	return info
}

func List() []VideoInfo {
	var info []VideoInfo
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		b.ForEach(func(_, v []byte) error {
			var i VideoInfo
			json.Unmarshal(v, &i)
			info = append(info, i)
			return nil
		})
		return nil
	})
	return info
}

func Put(info VideoInfo) {
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		val, _ := json.Marshal(info)
		b.Put([]byte(info.ID), val)
		return nil
	})
}
