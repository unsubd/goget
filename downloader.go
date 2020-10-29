package main

import (
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
	"os"
	"rip/computeutils"
	"rip/ioutils"
	"sort"
)

func download(url string) (int64, error) {
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
	for _, batch := range batches {
		go downloadBatch(url, batchSize, batch[0], batch[1], fmt.Sprintf("%s%s-%s", temp, fileName, uniqueId), ch)
	}

	var parts []string

	for i := 0; i < len(batches); i++ {
		parts = append(parts, <-ch)
	}

	sort.Strings(parts)

	for _, partName := range parts {
		ioutils.AppendToFile(partName, fileName, batchSize)
	}
	return 1, nil
}

func downloadBatch(url string, size int, start int64, end int64, fileName string, s chan string) (int64, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return -1, err
	}
	req.Header.Set("range", fmt.Sprintf("bytes=%d-%d", start, end))
	res, err := client.Do(req)
	if err != nil {
		return -1, err
	}

	defer res.Body.Close()
	bytes := make([]byte, size)
	read, err := res.Body.Read(bytes)
	if err != nil && err != io.EOF {
		return -1, err
	}

	bytes = bytes[:read]
	ioutils.WriteToFile(bytes, fileName)
	s <- fileName
	return int64(size), nil
}
