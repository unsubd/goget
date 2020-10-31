package ioutils

import (
	"fmt"
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

func PrintTrack(trackingChannel chan int64, uniqueId string, fileName string, contentLength int64) {
	for i := range trackingChannel {
		logging.ConsoleOut(fmt.Sprintf("DOWNLOAD_STATUS %s %s : %f Done", fileName, uniqueId, float64(i)*100/float64(contentLength)))
	}

	logging.ConsoleOut("\nDOWNLOAD_COMPLETE", fileName, uniqueId)
}