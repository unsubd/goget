package ioutils

import (
	"fmt"
	"goget/computeutils"
	"goget/constants"
	"goget/logging"
	"time"
)

func Track(uniqueId string, directoryPath string) (chan int64, chan bool) {
	ch := make(chan int64)
	stopChannel := make(chan bool)
	logging.LogDebug("TRACKING", uniqueId)
	go func() {
		for true {
			size, err := GetTotalFileSize(uniqueId, directoryPath)
			if err != nil {
				logging.LogError("GET_TOTAL_SIZE", err, uniqueId)
				return
			}
			select {
			case stop := <-stopChannel:
				if stop {
					logging.LogDebug("STOP_TRACKING", uniqueId)
					close(ch)
					close(stopChannel)
					return
				}
			case ch <- size:
			}
			time.Sleep(1 * time.Second)
		}
	}()

	return ch, stopChannel
}

func PrintTrack(trackingChannel chan constants.DownloadStatus) {
	lineNumbers := make(map[string]int)
	count := 0
	for i := range trackingChannel {
		var val int
		val, ok := lineNumbers[i.Id]
		if !ok {
			lineNumbers[i.Id] = count
			val = count
			count++
		}
		printStatus(i, val)
	}
	ConsoleOut("Download Complete", len(lineNumbers))
}

func printStatus(status constants.DownloadStatus, y int) {
	filePath := computeutils.GetFilePath(status.Dir, status.FileName)
	percentCompletion := (float64(status.Downloaded) * 100.0) / float64(status.Total)
	currentOperation := status.Op
	message := fmt.Sprintf("%s %0.2f %s", filePath, percentCompletion, currentOperation)
	ConsoleOut(message, y)
}
