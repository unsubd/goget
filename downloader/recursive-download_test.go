package downloader

import (
	"fmt"
	"goget/constants"
	"testing"
	"time"
)

func TestDownloadRecursive(t *testing.T) {
	ch, err := DownloadRecursive("http://localhost:8000", 100, "", 500*constants.MegaByte)
	if err != nil {
		t.Fail()
	}

	for status := range ch {
		fmt.Println(status)
		time.Sleep(1 * time.Second)
	}
}
