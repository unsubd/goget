package ioutils

import (
	"goget/logging"
	"net/http"
)

func RemoteFileSize(url string) (int64, error) {
	res, err := http.Head(url)

	if err != nil {
		logging.LogError("REMOTE_FILE_SIZE", err, url)
		return -1, err
	}
	return res.ContentLength, nil
}
