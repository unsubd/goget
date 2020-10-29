package main

import (
	"fmt"
	"log"
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
	fileName := computeutils.FileNameFromUrl(url)
	checksum, err := cryptoutils.FileChecksumSHA256(fileName)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Download complete: %s\n", fileName)
	fmt.Printf("SHA-256 checksum : %v\n", checksum)
}
