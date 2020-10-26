package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
)

func download(url string) (int64, error) {
	const batchSize = 10 * 1000 * 1000 // 10 MB
	contentLength, err := contentLength(url)
	if err != nil {
		return 0, err
	}
	if contentLength == -1 {
		contentLength = 9223372036854775807
	}

	batches := batch(contentLength, batchSize)
	fileName := extractFileName(url)

	for i, batch := range batches {
		downloadBatch(url, batchSize, batch[0], batch[1], fmt.Sprintf("%s-%d", fileName, i))
	}

	return 1, nil
}

func contentLength(url string) (int64, error) {
	res, err := http.Head(url)

	if err != nil {
		return -1, err
	}
	return res.ContentLength, nil
}
func downloadBatch(url string, size int, start int64, end int64, fileName string) (int64, error) {
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

	writeToFile(res.Body, fileName, size)

	return int64(size), nil
}

func writeToFile(body io.ReadCloser, fileName string, size int) {
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for true {
		bytes := make([]byte, size)
		read, err := body.Read(bytes)

		bytes = bytes[:read]

		writer.Write(bytes)
		if err != nil || read == 0 {
			break
		}
	}
}
