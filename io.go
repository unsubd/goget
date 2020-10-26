package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
)

func download(url string) (int64, error) {

	res, err := http.Get(url)

	if err != nil {
		return -1, err
	}
	defer res.Body.Close()
	writeToFile(res.Body, extractFileName(url), res.ContentLength)

	return res.ContentLength, nil
}

func writeToFile(body io.ReadCloser, fileName string, contentLength int64) {
	len := contentLength
	if len < 0 {
		len = 9223372036854775807
	}

	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()
	const limit = 10 * 1000 * 1000 // 10 MB

	for true {
		bytes := make([]byte, limit)
		read, err := body.Read(bytes)

		bytes = bytes[:read]

		writer.Write(bytes)
		if err != nil || read == 0 {
			break
		}
		len -= limit
	}
}
