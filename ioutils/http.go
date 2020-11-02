package ioutils

import (
	"goget/computeutils"
	"goget/logging"
	"io"
	"net/http"
	"net/url"
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

func GetDownloadLinks(baseUrl string) ([]string, error) {
	parsedUrl, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}

	body, err := HttpGet(parsedUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	a, err := computeutils.ParseHtml(body)
	if err != nil {
		return nil, err
	}

	fileNames := computeutils.ExtractFileNames(a)

	var results []string
	for _, fileName := range fileNames {
		parse, err := parsedUrl.Parse(fileName)
		if err == nil {
			results = append(results, parse.String())
		}
	}

	return results, nil
}
