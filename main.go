package main

import (
	"flag"
	"fmt"
	"goget/computeutils"
	"goget/constants"
	"goget/downloader"
	"goget/ioutils"
	"log"
	"os"
)

func main() {
	var size int64
	flag.Int64Var(&size, "m", 500, "MegaBytes of size to allocate for download")

	var url string
	flag.StringVar(&url, "url", "", "URL to download")

	var recursionDepth int
	flag.IntVar(&recursionDepth, "r", 1, "Recursion Depth")

	var outputDirectory string
	flag.StringVar(&outputDirectory, "o", ".", "Output Directory")

	flag.Parse()
	if url == "" {
		ioutils.ConsoleOutLn("URL CANNOT BE EMPTY")
		log.Fatal("URL CANNOT BE EMPTY")
	}

	logDir := computeutils.GetFilePath(outputDirectory, "logs")
	err3 := os.MkdirAll(logDir, 0777)
	if err3 != nil {
		log.Println(err3)
	}
	logFile, err2 := os.OpenFile(computeutils.GetFilePath(logDir, "app.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err2 != nil {
		log.Fatal(err2)
	}

	log.SetOutput(logFile)

	statusChannel, err := downloader.DownloadRecursive(url, recursionDepth, outputDirectory, size*constants.MegaByte)

	ioutils.PrintTrack(statusChannel)

	if err != nil {
		ioutils.ConsoleOutLn(fmt.Sprintf("ERROR WHILE DOWNLOADING FILE %v", err))
	}
}
