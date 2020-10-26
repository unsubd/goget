package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
)

func download(url string) (int64, error) {
	size := 10 * 1000 * 1000 // 10 MB
	return downloadBatch(url, size, 0, 9223372036854775807, extractFileName(url), 0)
}

func downloadBatch(url string, size int, start int64, end int64, fileName string, index int) (int64, error) {
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
