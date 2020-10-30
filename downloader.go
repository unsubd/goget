package main

import (
	"fmt"
	"github.com/google/uuid"
	"goget/computeutils"
	"goget/ioutils"
	"goget/logging"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func downloadFile(url string) (int64, error) {
	const batchSize = 10 * 1000 * 1000 // 10 MB
	logging.LogDebug("DOWNLOAD STARTING FOR", url)
	contentLength, err := ioutils.RemoteFileSize(url)
	logging.LogDebug("CONTENT_LENGTH", fmt.Sprintf("bytes %v", contentLength), fmt.Sprintf("mb %v", contentLength/(1000*1000)))
	if err != nil {
		return -1, err
	}

	batches := computeutils.CreateBatches(contentLength, batchSize)
	logging.LogDebug("BATCH_COUNT", len(batches), url)
	fileName := computeutils.FileNameFromUrl(url)
	logging.LogDebug("FILE_NAME", fileName, url)
	ch := make(chan string, len(batches))
	temp := os.TempDir()
	uniqueId := uuid.New().String()
	logging.LogDebug("UUID", uniqueId, url)
	baseFileName := fmt.Sprintf("%s%s-%s", temp, fileName, uniqueId)
	for i, batch := range batches {
		start := batch[0]
		end := batch[1]
		index := i
		go func() {
			filePartName := fmt.Sprintf("%s-%d", baseFileName, index)
			err2 := downloadFilePart(url, start, end, filePartName)
			if err2 != nil {
				logging.LogError("DOWNLOADING_PART", err, filePartName)
				ch <- fmt.Sprintf("%s-ERROR-%s", fileName, err2.Error())
			} else {
				ch <- filePartName
			}

		}()
	}
	trackingChannel := computeutils.Track(uniqueId, temp)

	go func() {
		for i := range trackingChannel {
			logging.LogDebug(fmt.Sprintf("DOWNLOAD_STATUS %s %s : %f Done\n", fileName, uniqueId, float64(i)*100/float64(contentLength)))
		}
	}()

	for i := 0; i < len(batches); i++ {
		part := <-ch
		if strings.Contains(part, "ERROR") {
			trackingChannel <- 1
			break
		}
	}

	for i := 0; i < len(batches); i++ {
		err := ioutils.AppendToFile(fmt.Sprintf("%s-%d", baseFileName, i), fileName, batchSize)
		if err != nil {
			logging.LogError("APPEND_TO_FILE", err, fileName, baseFileName)
			log.Println(err.Error())
		}
	}

	logging.LogDebug("DOWNLOAD_COMPLETE", url, uniqueId)
	return 1, nil
}

func downloadFilePart(url string, start int64, end int64, fileName string) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		logging.LogError("HTTP_REQUEST", err, fileName, url)
		return err
	}
	req.Header.Set("range", fmt.Sprintf("bytes=%d-%d", start, end))
	res, err := client.Do(req)
	if err != nil {
		logging.LogError("HTTP_GET", err, fileName, url)
		return err
	}

	defer res.Body.Close()
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil && err != io.EOF {
		logging.LogError("HTTP_GET_BODY", err, fileName, url)
		return err
	}

	ioutils.WriteToFile(bytes, fileName)
	return nil
}
