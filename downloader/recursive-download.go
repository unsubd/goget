package downloader

import (
	"fmt"
	"github.com/google/uuid"
	"goget/computeutils"
	"goget/constants"
	"goget/ioutils"
	"goget/logging"
	"os"
	"strings"
)

func DownloadRecursive(url string, depth int, directory string, limit constants.Size) (chan constants.DownloadStatus, error) {
	ch := make(chan constants.DownloadStatus)

	go func() {
		var rec func(string, int, string)
		defer close(ch)
		baseDirectory := directory
		tempDirectory := ioutils.GetTempDir()
		downloadId := uuid.New().String()

		if baseDirectory != "." {
			tempDirectory = computeutils.GetFilePath(baseDirectory, fmt.Sprintf("temp%s", downloadId))
			err := os.MkdirAll(tempDirectory, 0777)
			if err != nil {
				logging.LogError(err)
				tempDirectory = ioutils.GetTempDir()
			}
		}
		logging.LogDebug("Temp Directory", tempDirectory, downloadId)

		rec = func(url string, depth int, baseDirectory string) {

			if depth > 0 {
				links, _ := ioutils.GetDownloadLinks(url)
				directories := processLinks(links, limit, ch, baseDirectory, tempDirectory)
				depth--
				for _, directoryUrl := range directories {
					rec(directoryUrl, depth, computeutils.GetFilePath(baseDirectory, computeutils.FileNameFromUrl(url)))
				}
			}

		}

		rec(url, depth, baseDirectory)

	}()

	return ch, nil
}

func processLinks(links []string, limit constants.Size, response chan constants.DownloadStatus, dir string, tempDirectory string) []string {
	var directories []string
	for _, link := range links {
		_, contentType, _, err := ioutils.HttpRequest("HEAD", link, nil, nil)
		if err != nil {
			response <- constants.DownloadStatus{Error: err}
		}

		if !strings.Contains(contentType, "text/html") {
			trackingChannel, uniqueId, contentLength, fileName, err := Download(link, limit, dir, tempDirectory)
			for s := range trackingChannel {
				response <- constants.DownloadStatus{Downloaded: s.downloaded,
					Id:       uniqueId,
					Total:    contentLength,
					FileName: fileName,
					Error:    err,
					Dir:      dir,
					Op:       s.op,
				}

			}
		} else {
			directories = append(directories, link)
		}
	}

	return directories
}
