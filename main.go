package main

import (
	"flag"
	"fmt"
	"goget/constants"
	"goget/downloader"
	"goget/ioutils"
	"goget/logging"
	"log"
	"os"
)

func main() {
	err3 := os.MkdirAll("logs", 0777)
	if err3 != nil {
		log.Println(err3)
	}
	logFile, err2 := os.OpenFile("logs/app.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err2 != nil {
		log.Fatal(err2)
	}

	log.SetOutput(logFile)
	var size int64
	flag.Int64Var(&size, "m", 500, "MegaBytes of size to allocate for download")

	var url string
	flag.StringVar(&url, "url", "", "URL to download")

	var recursionDepth int
	flag.IntVar(&recursionDepth, "r", 1, "Recursion Depth")

	flag.Parse()
	if url == "" {
		logging.ConsoleOut("URL CANNOT BE EMPTY")
		log.Fatal("URL CANNOT BE EMPTY")
	}

	statusChannel, err := downloader.DownloadRecursive(url, recursionDepth, size*constants.MegaByte)

	if err3 != nil {

	}

	ioutils.PrintTrack(statusChannel)

	if err != nil {
		logging.ConsoleOut(fmt.Sprintf("MAIN ERROR DOWNLOADING FILE: %v", err))
	}
}
