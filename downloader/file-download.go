package downloader

import (
	"fmt"
	"github.com/google/uuid"
	"goget/computeutils"
	"goget/constants"
	"goget/ioutils"
	"goget/logging"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync/atomic"
)

func Download(url string, limit constants.Size, dir string, temp string, resume bool) (chan struct {
	downloaded int64
	op         string
}, string, int64, string, error) {
	fileDownloadTracker := make(chan struct {
		downloaded int64
		op         string
	})
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
	uniqueId := uuid.New().String()

	logging.LogDebug("UUID", uniqueId, url)
	skips := make(map[int]bool)
	id := ""
	if resume {
		fileNames, _ := ioutils.GetFilesFromPattern(fileName, temp)
		skips, id = computeutils.ExtractResumeMetaData(fileNames, fileName, temp, batchSize)
		if id != "" {
			uniqueId = id
			fmt.Println("Resuming", uniqueId, len(skips))
		}
	}
	baseFileName := fmt.Sprintf("%s-%s", computeutils.GetFilePath(temp, fileName), uniqueId)
	_, err = os.Stat(computeutils.GetFilePath(dir, fileName))
	filePresent := true
	if stat, err := os.Stat(computeutils.GetFilePath(dir, fileName)); os.IsNotExist(err) {
		filePresent = false
	} else {
		filePresent = contentLength == stat.Size()
		if !filePresent {
			ioutils.DeleteFiles(computeutils.GetFilePath(dir, fileName))
		}
	}

	if filePresent && resume {
		defer close(fileDownloadTracker)
		return fileDownloadTracker, uniqueId, contentLength, fileName, os.ErrExist
	}

	go dispatchBatches(url, batches, baseFileName, ch, skips, fileName, uniqueId, int(limit/batchSize))
	trackingChannel, stopChannel := ioutils.Track(uniqueId, temp)

	go func() {
		downloaded := int64(0)

		for i := range trackingChannel {
			downloaded = i
			fileDownloadTracker <- struct {
				downloaded int64
				op         string
			}{downloaded: i, op: "DOWNLOADING"}
		}

		fileDownloadTracker <- struct {
			downloaded int64
			op         string
		}{downloaded: downloaded, op: "APPENDING"}

	}()
	go func() {
		defer ioutils.DeleteFiles(baseFileName)
		defer close(fileDownloadTracker)
		for part := range ch {
			if strings.Contains(part, "ERROR") {
				break
			}
		}

		stopChannel <- true
		for i := 0; i < len(batches); i++ {
			err := ioutils.AppendToFile(fmt.Sprintf("%s-%d", baseFileName, i), fileName, dir, batchSize)
			if err != nil {
				logging.LogError("APPEND_TO_FILE", err, fileName, baseFileName)
			}
		}

		fileDownloadTracker <- struct {
			downloaded int64
			op         string
		}{downloaded: contentLength, op: "DONE"}

		logging.LogDebug("DOWNLOAD_COMPLETE", url, uniqueId)
	}()

	return fileDownloadTracker, uniqueId, contentLength, fileName, nil
}

func dispatchBatches(url string, batches [][]int64, baseFileName string, response chan string, skips map[int]bool, fileName string, uniqueId string, limit int) {
	dispatchChannel := make(chan string)
	logging.LogDebug("DISPATCH", fmt.Sprintf("LIMIT=%v", limit), uniqueId)
	count := 0
	start := 0
	remaining := int32(len(batches)) - int32(len(skips))

	for index := 0; count < limit && index < len(batches); index++ {
		_, ok := skips[index]
		if !ok {
			dispatch(url, baseFileName, index, batches[index][0], batches[index][1], dispatchChannel, fileName)
			skips[index] = true
			count++
		}
		start = index
	}

	index := int32(start + 1)
	logging.LogDebug("DISPATCH_PARTIAL", fmt.Sprintf("START=%v", index), uniqueId)

	go func() {
		defer close(response)
		defer close(dispatchChannel)
	loop:
		for remaining > 0 {
			filePart := <-dispatchChannel
			response <- filePart
			atomic.AddInt32(&remaining, -1)
			if remaining <= 0 || strings.Contains(filePart, "ERROR") {
				logging.LogDebug("DISPATCH_DONE", uniqueId)
				break loop
			}

			if index < int32(len(batches)) {
				logging.LogDebug("DISPATCH_PARTIAL", index, uniqueId)
				_, ok := skips[int(index)]
				if !ok {
					dispatch(url, baseFileName, int(index), batches[index][0], batches[index][1], dispatchChannel, fileName)
				}
			}

			atomic.AddInt32(&index, 1)
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
	body, _, _, err := ioutils.HttpRequest("GET", url, map[string]string{"range": fmt.Sprintf("bytes=%d-%d", start, end)}, nil)
	if err != nil {
		logging.LogError("HTTP_GET", err, fileName, url)
		return err
	}
	defer body.Close()
	bytes, err := ioutil.ReadAll(body)
	if err != nil && err != io.EOF {
		logging.LogError("HTTP_GET_BODY", err, fileName, url)
		return err
	}

	ioutils.WriteToFile(bytes, fileName, "")
	return nil
}
