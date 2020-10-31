# goget
Download stuff from the internet with ease.

The idea is to have one command line tool to help me download an entire index of files, if need be.

Usage:

`$ go build goget`

`$ ./goget -m=100 -url=https://raw.githubusercontent.com/unsubd/geektrust-family/master/input.txt`
`This would run goget with 100 MB reserved for downloading objects`

TODO:

  1. Download Directories and sub directories
  2. Add a Tracker for downloads
