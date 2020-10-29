package computeutils

import (
	"strings"
)

func FileNameFromUrl(url string) string {
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
		start += size
		end = start + size
		if start > limit {
			break
		}
	}

	return batches
}
