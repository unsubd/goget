package ioutils

import (
	"bufio"
	"fmt"
	"goget/logging"
	"io/ioutil"
	"os"
	"strings"
)

func WriteToFile(bytes []byte, fileName string) {
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logging.LogError("WRITE_TO_FILE", err, fileName)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()
	writer.Write(bytes)
}

func AppendToFile(partialFilePath string, finalFilePath string, size int) error {
	file, err := os.Open(partialFilePath)
	logging.LogDebug("APPEND_TO_FILE CALLED", partialFilePath)
	if err != nil {
		logging.LogError("APPEND_TO_FILE", err, partialFilePath)
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

	logging.LogDebug("APPEND_TO_FILE SUCCESSFUL", partialFilePath)

	return nil
}

func GetTotalFileSize(pattern string, directoryPath string) (int64, error) {
	logging.LogDebug("GET_TOTAL_SIZE", pattern)
	var size int64
	files, err := ioutil.ReadDir(directoryPath)

	if err != nil {
		logging.LogError("GET_TOTAL_SIZE", err, pattern)
		return -1, err
	}

	var fileNames []string

	for _, file := range files {
		fileName := file.Name()
		if strings.Contains(fileName, pattern) {
			fileNames = append(fileNames, fileName)
		}
	}

	logging.LogDebug("FILE_COUNT", pattern, len(fileNames))

	for _, fileName := range fileNames {
		stat, err := os.Stat(fmt.Sprintf("%s%s", directoryPath, fileName))
		if err != nil {
			logging.LogError("OS_STAT", err, fileName)
			return -1, err
		}
		size += stat.Size()
	}

	logging.LogDebug("GET_TOTAL_SIZE_SUCCESSFUL", pattern, size)

	return size, nil
}
