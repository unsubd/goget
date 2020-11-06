package computeutils

import (
	"fmt"
	"goget/constants"
	"os"
	"regexp"
	"strconv"
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
	if dir[len(dir)-1] == '/' || dir[len(dir)-1] == '\\' {
		dir = dir[0 : len(dir)-1]
	}

	if fileName[0] == '/' || fileName[0] == '\\' {
		fileName = fileName[1:]
	}
	return fmt.Sprintf("%s/%s", dir, fileName)
}

func ExtractResumeMetaData(fileNames []string, name string, dir string, minSize constants.Size) (map[int]bool, string) {
	skips := make(map[int]bool)
	uniqueId := ""
	regex := regexp.MustCompile("^.*-(\\d+)$")

	for _, fileName := range fileNames {
		submatch := regex.FindStringSubmatch(fileName)
		atoi, _ := strconv.Atoi(submatch[1])
		stat, err := os.Stat(GetFilePath(dir, fileName))
		if err == nil && stat.Size() > minSize {
			skips[atoi] = true
		} else {
			os.Remove(GetFilePath(dir, fileName))
		}
		if uniqueId == "" {
			uniqueId = regexp.MustCompile(fmt.Sprintf("^.*%s-(.*)-\\d+$", name)).FindStringSubmatch(fileName)[1]
		}
	}

	return skips, uniqueId

}
