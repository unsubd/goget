package ioutils

import (
	"fmt"
	"goget/computeutils"
	"goget/logging"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func RemoteFileSize(url string) (int64, error) {
	res, err := http.Head(url)

	if err != nil {
		logging.LogError("REMOTE_FILE_SIZE", err, url)
		return -1, err
	}
	return res.ContentLength, nil
}

func HttpRequest(method string, url string, headers map[string]string, body io.Reader) (io.ReadCloser, string, int64, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, "", 0, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	res, err := client.Do(req)

	if err != nil {
		return nil, "", 0, err
	}

	return res.Body, res.Header.Get("content-type"), res.ContentLength, nil
}

func GetDownloadLinks(baseUrl string) ([]string, error) {
	fmt.Println("called")
	parsedUrl, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}

	body, contentType, _, err := HttpRequest("HEAD", parsedUrl.String(), nil, nil)
	if err != nil {
		return nil, err
	}

	if !strings.Contains(contentType, "text/html") {
		return []string{baseUrl}, nil
	}

	body, _, _, err = HttpRequest("GET", parsedUrl.String(), nil, nil)

	if err != nil {
		return nil, err
	}

	htmlRootNode, err := computeutils.ParseHtml(body)
	if err != nil {
		return nil, err
	}

	fileNames := computeutils.ExtractFileNames(htmlRootNode)

	var results []string
	for _, fileName := range fileNames {
		parse, err := parsedUrl.Parse(fileName)
		if err == nil {
			results = append(results, parse.String())
		}
	}

	return results, nil
}
