package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
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
	ch := make(chan string, len(batches))

	for i, batch := range batches {
		go downloadBatch(url, batchSize, batch[0], batch[1], fmt.Sprintf("%s-%d", fileName, i), ch)
	}

	var parts []string

	for i := 0; i < len(batches); i++ {
		parts = append(parts, <-ch)
	}

	sort.Strings(parts)

	for _, partName := range parts {
		merge(partName, fileName, batchSize)
	}
	return 1, nil
}

func merge(partName string, destinationName string, size int) error {
	file, err := os.Open(partName)
	if err != nil {
		return err
	}

	defer os.Remove(partName)

	reader := bufio.NewReader(file)

	for true {
		bytes := make([]byte, size)
		read, err := reader.Read(bytes)
		bytes = bytes[:read]
		writeToFile(bytes, destinationName)
		if err != nil || read == 0 {
			break
		}
	}

	return nil

}

func contentLength(url string) (int64, error) {
	res, err := http.Head(url)

	if err != nil {
		return -1, err
	}
	return res.ContentLength, nil
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
	writeToFile(bytes, fileName)
	s <- fileName
	return int64(size), nil
}

func writeToFile(bytes []byte, fileName string) {
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()
	writer.Write(bytes)
}
