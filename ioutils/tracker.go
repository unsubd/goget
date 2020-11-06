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

func PrintTrack(trackingChannel chan constants.DownloadStatus, doneCallback func(status constants.DownloadStatus)) {
	downloaded := make(map[string]bool)
	for i := range trackingChannel {
		_, ok := downloaded[i.Id]
		if !ok {
			NewLine()
			downloaded[i.Id] = true
		} else {
			goToStart()
		}
		printStatus(i)
		if i.Op == "DONE" {
			doneCallback(i)
		}
	}
	NewLine()
	ConsoleOut("Download Complete!")
}

func printStatus(status constants.DownloadStatus) {
	filePath := computeutils.GetFilePath(status.Dir, status.FileName)
	percentCompletion := (float64(status.Downloaded) * 100.0) / float64(status.Total)
	currentOperation := status.Op
	message := fmt.Sprintf("%s %0.2f %d/%d %s", filePath, percentCompletion, status.Downloaded, status.Total, currentOperation)
	ConsoleOut(message)
}

func loop(message string) {
	symbols := []string{"|", "/", "|", "\\"}
	count := 5
	for i := 0; i < count; i++ {
		for j := 0; j < len(symbols); j++ {
			ConsoleOut(fmt.Sprintf("%s  %s", message, symbols[j]))
			time.Sleep(50 * time.Millisecond)
		}
	}
}
