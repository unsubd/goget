package computeutils

import (
	"fmt"
	"goget/ioutils"
	"time"
)

func Track(uniqueId string, directoryPath string) chan int64 {
	ch := make(chan int64)

	go func() {
		for true {
			size, err := ioutils.GetTotalFileSize(uniqueId, directoryPath)
			if err != nil {
				fmt.Println(err)
				return
			}
			select {
			case stop := <-ch:
				if stop == 1 {
					close(ch)
					return
				}
			case ch <- size:
			}
			time.Sleep(1 * time.Second)
		}
	}()

	return ch
}
