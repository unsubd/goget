package computeutils

import (
	"fmt"
	"strings"
)

func FileNameFromUrl(url string) string {
	if lastIndex := strings.LastIndex(url, "/"); lastIndex == len(url)-1 {
		url = url[:lastIndex]
	}
	return url[strings.LastIndex(url, "/")+1:]
}

func CreateBatches(limit int64, size int64) [][]int64 {
	if size >= limit {
		return [][]int64{{0, limit}}
	}

	var batches [][]int64

	start := int64(0)
	end := size

	for true {
		batches = append(batches, []int64{start, end})
		start += size + 1
		end = start + size
		if start > limit {
			break
		}
	}

	return batches
}

func GetFilePath(dir string, fileName string) string {
	if dir[len(dir)-1] == '/'|'\\' {
		dir = dir[0 : len(dir)-1]
	}
	return fmt.Sprintf("%s/%s", dir, fileName)
}
