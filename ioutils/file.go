package ioutils

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func WriteToFile(bytes []byte, fileName string) {
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()
	writer.Write(bytes)
}

func AppendToFile(partialFilePath string, finalFilePath string, size int) error {
	file, err := os.Open(partialFilePath)
	if err != nil {
		return err
	}

	defer os.Remove(partialFilePath)

	reader := bufio.NewReader(file)

	for true {
		bytes := make([]byte, size)
		read, err := reader.Read(bytes)
		bytes = bytes[:read]
		WriteToFile(bytes, finalFilePath)
		if err != nil || read == 0 {
			break
		}
	}

	return nil
}

func GetTotalFileSize(pattern string, directoryPath string) (int64, error) {
	var size int64
	files, err := ioutil.ReadDir(directoryPath)

	if err != nil {
		return -1, err
	}

	var fileNames []string

	for _, file := range files {
		fileName := file.Name()
		if strings.Contains(fileName, pattern) {
			fileNames = append(fileNames, fileName)
		}
	}

	for _, fileName := range fileNames {
		stat, err := os.Stat(fmt.Sprintf("%s%s", directoryPath, fileName))
		if err != nil {
			return -1, err
		}
		size += stat.Size()
	}

	return size, nil
}
