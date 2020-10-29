package main

import (
	"fmt"
	"github.com/google/uuid"
	"goget/computeutils"
	"goget/ioutils"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func downloadFile(url string) (int64, error) {
	const batchSize = 10 * 1000 * 1000 // 10 MB
	contentLength, err := ioutils.RemoteFileSize(url)
	if err != nil {
		return 0, err
	}
	if contentLength == -1 {
		contentLength = 9223372036854775807
	}

	batches := computeutils.CreateBatches(contentLength, batchSize)
	fileName := computeutils.FileNameFromUrl(url)
	ch := make(chan string, len(batches))
	temp := os.TempDir()
	uniqueId := uuid.New().String()
	baseFileName := fmt.Sprintf("%s%s-%s", temp, fileName, uniqueId)
	for i, batch := range batches {
		start := batch[0]
		end := batch[1]
		index := i
		go func() {
			filePartName := fmt.Sprintf("%s-%d", baseFileName, index)
			err2 := downloadFilePart(url, start, end, filePartName)
			if err2 != nil {
				ch <- fmt.Sprintf("%s-ERROR-%s", fileName, err2.Error())
			} else {
				ch <- filePartName
			}

		}()
	}
	trackingChannel := computeutils.Track(uniqueId, temp)

	go func() {
		for i := range trackingChannel {
			fmt.Printf("Download Status %s : %f Done\n", fileName, float64(i)*100/float64(contentLength))
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
			fmt.Println(err.Error())
		}
	}
	return 1, nil
}

func downloadFilePart(url string, start int64, end int64, fileName string) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return err
	}
	req.Header.Set("range", fmt.Sprintf("bytes=%d-%d", start, end))
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil && err != io.EOF {
		return err
	}

	ioutils.WriteToFile(bytes, fileName)
	return nil
}
