package main

import "strings"

func extractFileName(url string) string {
	return url[strings.LastIndex(url, "/")+1:]
}
