package main

import (
	"fmt"
	"os"
)

func main() {
	url := os.Args[1]
	_, err := download(url)
	if err != nil {
		fmt.Println(err.Error())
	}
	checksum, err := checksum(extractFileName(url))
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(checksum)
}
