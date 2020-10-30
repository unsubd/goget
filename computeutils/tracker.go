package computeutils

import (
	"goget/ioutils"
	"goget/logging"
	"time"
)

func Track(uniqueId string, directoryPath string) (chan int64, chan bool) {
	ch := make(chan int64)
	stopChannel := make(chan bool)
	logging.LogDebug("TRACKING", uniqueId)
	go func() {
		for true {
			size, err := ioutils.GetTotalFileSize(uniqueId, directoryPath)
			if err != nil {
				logging.LogError("GET_TOTAL_SIZE", err, uniqueId)
				return
			}
			select {
			case stop := <-stopChannel:
				if stop {
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
