package downloader

import (
	"fmt"
	"goget/computeutils"
	"goget/constants"
	"goget/ioutils"
	"strings"
)

func DownloadRecursive(url string, depth int, limit constants.Size) (chan constants.Status, error) {
	ch := make(chan constants.Status)

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

func processLinks(links []string, limit constants.Size, response chan constants.Status, dir string) []string {
	var directories []string
	for _, link := range links {
		_, contentType, _, err := ioutils.HttpRequest("HEAD", link, nil, nil)
		if err != nil {
			response <- constants.Status{Error: err}
		}

		if !strings.Contains(contentType, "text/html") {
			trackingChannel, uniqueId, contentLength, fileName, err := Download(link, limit, dir)
			for downloadedSize := range trackingChannel {
				response <- constants.Status{Downloaded: downloadedSize,
					Id:       uniqueId,
					Total:    contentLength,
					FileName: fileName,
					Error:    err,
				}
			}
		} else {
			directories = append(directories, link)
		}
	}

	return directories
}
