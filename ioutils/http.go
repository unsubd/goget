package ioutils

import (
	"goget/logging"
	"io"
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

func HttpGet(url string, headers map[string]string) (io.ReadCloser, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	return res.Body, nil
}
