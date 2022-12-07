# Sync

sync videos from other platform to VolcEngine-VOD

## Development

setting your credentials in env

```
export VOD_AK=xxx
export VOD_SK=xxx
export VOD_SPACE=xxx
```

start the backend server

```
go run $(ls -1 *.go | grep -v _test.go)
```

run frontend

```
cd web
yarn start
```

## How to use the app

1. click `Upload Video` to provide a list of source_url
2. waiting the backend to send sync request to VolcEngine, the list of ids will refresh
3. click on any id to watch the video (with source_url)
4. the video panel will show if it's a source url or volc video
5. when the reason is "not published", go to VolcEngine to publish your video
6. refresh the page, you can view the same video with after-sync VolcEngine url

## References - Getting SourceURI

Tencent Cloud

1. Go to https://console.cloud.tencent.com/vod/media
2. Select All files
3. Batch Process -> Export Files Information
4. In the csv file downloaded, you will find all source_url
