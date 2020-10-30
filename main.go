package main

import (
	"fmt"
	"goget/computeutils"
	"goget/cryptoutils"
	"log"
	"os"
)

func main() {
	url := os.Args[1]
	_, err := downloadFile(url)
	if err != nil {
		fmt.Printf("MAIN ERROR DOWNLOADING FILE: %v", err)
	}
	fileName := computeutils.FileNameFromUrl(url)
	checksum, err := cryptoutils.FileChecksumSHA256(fileName)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Download complete: %s\n", fileName)
	fmt.Printf("SHA-256 checksum : %v\n", checksum)
}
