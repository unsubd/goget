package downloader

import (
	"fmt"
	"goget/computeutils"
	"goget/constants"
	"goget/ioutils"
	"strings"
)

type Status struct {
	id         string
	downloaded int64
	total      int64
	fileName   string
	error      error
}

func DownloadRecursive(url string, depth int, limit constants.Size) (chan Status, error) {
	ch := make(chan Status)

	go func() {
		var rec func(string, int, string)
		defer close(ch)
		baseDirectory := "."

		rec = func(url string, depth int, baseDirectory string) {

			if depth > 0 {
				links, _ := ioutils.GetDownloadLinks(url)
				directories := processLinks(links, limit, ch, baseDirectory)
				depth--
				for _, directoryUrl := range directories {
					rec(directoryUrl, depth, fmt.Sprintf("%s/%s", baseDirectory, computeutils.FileNameFromUrl(directoryUrl)))
				}
			}

		}

		rec(url, depth, baseDirectory)

	}()

	return ch, nil
}

func processLinks(links []string, limit constants.Size, response chan Status, dir string) []string {
	var directories []string
	for _, link := range links {
		_, contentType, _, err := ioutils.HttpRequest("HEAD", link, nil, nil)
		if err != nil {
			response <- Status{error: err}
		}

		if !strings.Contains(contentType, "text/html") {
			trackingChannel, uniqueId, contentLength, fileName, err := Download(link, limit, dir)
			for downloadedSize := range trackingChannel {
				response <- Status{downloaded: downloadedSize,
					id:       uniqueId,
					total:    contentLength,
					fileName: fileName,
					error:    err,
				}
			}
		} else {
			directories = append(directories, link)
		}
	}

	return directories
}
