package main

import (
	"fmt"
	"os"
	"rip/computeutils"
	"rip/cryptoutils"
)

func main() {
	url := os.Args[1]
	_, err := download(url)
	if err != nil {
		fmt.Println(err.Error())
	}
	checksum, err := cryptoutils.FileChecksumSHA256(computeutils.FileNameFromUrl(url))
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(checksum)
}
