package ioutils

import "net/http"

func RemoteFileSize(url string) (int64, error) {
	res, err := http.Head(url)

	if err != nil {
		return -1, err
	}
	return res.ContentLength, nil
}
