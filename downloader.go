package main

import (
	"fmt"
	"github.com/google/uuid"
	"goget/computeutils"
	"goget/constants"
	"goget/ioutils"
	"goget/logging"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync/atomic"
)

func DownloadFile(url string, limit constants.Size) (chan int64, string, int64, string, error) {
	const batchSize = 10 * constants.MegaByte // 10 MB
	logging.LogDebug("DOWNLOAD STARTING FOR", url)
	contentLength, err := ioutils.RemoteFileSize(url)
	logging.LogDebug("CONTENT_LENGTH", fmt.Sprintf("bytes %v", contentLength), fmt.Sprintf("mb %v", contentLength/(1000*1000)))
	if err != nil {
		return nil, "", -1, "", err
	}

	batches := computeutils.CreateBatches(contentLength, batchSize)
	logging.LogDebug("BATCH_COUNT", len(batches), url)
	fileName := computeutils.FileNameFromUrl(url)
	logging.LogDebug("FILE_NAME", fileName, url)

	ch := make(chan string, len(batches))
	temp := ioutils.GetTempDir()
	uniqueId := uuid.New().String()

	logging.LogDebug("UUID", uniqueId, url)
	logging.ConsoleOut("UUID", uniqueId)
	baseFileName := fmt.Sprintf("%s%s-%s", temp, fileName, uniqueId)

	dispatchBatches(url, batches, baseFileName, ch, fileName, uniqueId, int(limit/batchSize))
	trackingChannel, stopChannel := ioutils.Track(uniqueId, temp)

	go func() {
		defer ioutils.DeleteFiles(baseFileName)
		for part := range ch {
			if strings.Contains(part, "ERROR") {
				break
			}
		}

		for i := 0; i < len(batches); i++ {
			err := ioutils.AppendToFile(fmt.Sprintf("%s-%d", baseFileName, i), fileName, batchSize)
			if err != nil {
				logging.LogError("APPEND_TO_FILE", err, fileName, baseFileName)
			}
		}

		stopChannel <- true
		logging.LogDebug("DOWNLOAD_COMPLETE", url, uniqueId)
	}()

	return trackingChannel, uniqueId, contentLength, fileName, nil
}

func dispatchBatches(url string, batches [][]int64, baseFileName string, response chan string, fileName string, uniqueId string, limit int) {
	dispatchChannel := make(chan string)
	logging.LogDebug("DISPATCH", fmt.Sprintf("LIMIT=%v", limit), uniqueId)
	for count := 0; count < limit; count++ {
		dispatch(url, baseFileName, count, batches[count][0], batches[count][1], dispatchChannel, fileName)
	}

	index := int32(limit)
	logging.LogDebug("DISPATCH_PARTIAL", fmt.Sprintf("START=%v", index), uniqueId)
	receiveCount := int32(0)

	go func() {
		defer close(response)
		defer close(dispatchChannel)
	loop:
		for true {
			select {
			case filePart := <-dispatchChannel:
				response <- filePart
				if index < int32(len(batches)) {
					logging.LogDebug("DISPATCH_PARTIAL", index, uniqueId)
					dispatch(url, baseFileName, int(index), batches[index][0], batches[index][1], dispatchChannel, fileName)
				}
				atomic.AddInt32(&index, 1)
				atomic.AddInt32(&receiveCount, 1)
				if receiveCount >= int32(len(batches)) || strings.Contains(filePart, "ERROR") {
					logging.LogDebug("DISPATCH_DONE", uniqueId)
					break loop
				}
			}
		}

	}()

}

func dispatch(url string, baseFileName string, index int, start int64, end int64, response chan string, fileName string) {
	go func() {
		logging.LogDebug("DISPATCHING", index, baseFileName)
		filePartName := fmt.Sprintf("%s-%d", baseFileName, index)
		err := downloadPartialFile(url, start, end, filePartName)
		if err != nil {
			logging.LogError("DOWNLOADING_PART", err, filePartName)
			response <- fmt.Sprintf("%s-ERROR-%s", fileName, err.Error())
		} else {
			response <- filePartName
		}

	}()
}

func downloadPartialFile(url string, start int64, end int64, fileName string) error {
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
