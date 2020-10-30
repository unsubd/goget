package main

import (
	"fmt"
	"goget/computeutils"
	"goget/constants"
	"goget/cryptoutils"
	"log"
	"os"
)

func main() {
	err3 := os.MkdirAll("/var/log/goget", 0644)
	if err3 != nil {
		log.Println(err3)
	}
	logFile, err2 := os.OpenFile("/var/log/goget/app.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err2 != nil {
		log.Fatal(err2)
	}

	log.SetOutput(logFile)

	url := os.Args[1]
	_, err := downloadFile(url, 500*constants.MegaByte)
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
