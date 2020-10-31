# goget
Download stuff from the internet with ease.

The idea is to have one command line tool to help me download an entire index of files, if need be.

Usage:

`$ go build goget`

`$ ./goget -m=100 -url=https://stuff.mit.edu/afs/sipb/contrib/pi/pi-billion.txt`

```
[UUID 91273edf-e095-4ff4-8c3c-128523698bd6]
[DOWNLOAD_STATUS pi-billion.txt 91273edf-e095-4ff4-8c3c-128523698bd6 : 0.000000 Done]
[DOWNLOAD_STATUS pi-billion.txt 91273edf-e095-4ff4-8c3c-128523698bd6 : 0.000000 Done]
[DOWNLOAD_STATUS pi-billion.txt 91273edf-e095-4ff4-8c3c-128523698bd6 : 0.000000 Done]
[DOWNLOAD_STATUS pi-billion.txt 91273edf-e095-4ff4-8c3c-128523698bd6 : 0.000000 Done]
[DOWNLOAD_STATUS pi-billion.txt 91273edf-e095-4ff4-8c3c-128523698bd6 : 0.000000 Done]
[DOWNLOAD_STATUS pi-billion.txt 91273edf-e095-4ff4-8c3c-128523698bd6 : 0.000000 Done]
[DOWNLOAD_STATUS pi-billion.txt 91273edf-e095-4ff4-8c3c-128523698bd6 : 0.000000 Done]
[DOWNLOAD_STATUS pi-billion.txt 91273edf-e095-4ff4-8c3c-128523698bd6 : 0.000000 Done]
.
.
.
[DOWNLOAD_STATUS pi-billion.txt 91273edf-e095-4ff4-8c3c-128523698bd6 : 83.000008 Done]
[DOWNLOAD_STATUS pi-billion.txt 91273edf-e095-4ff4-8c3c-128523698bd6 : 91.000009 Done]
[DOWNLOAD_STATUS pi-billion.txt 91273edf-e095-4ff4-8c3c-128523698bd6 : 92.000009 Done]
[DOWNLOAD_STATUS pi-billion.txt 91273edf-e095-4ff4-8c3c-128523698bd6 : 96.000009 Done]
[DOWNLOAD_STATUS pi-billion.txt 91273edf-e095-4ff4-8c3c-128523698bd6 : 97.000010 Done]
[DOWNLOAD_STATUS pi-billion.txt 91273edf-e095-4ff4-8c3c-128523698bd6 : 99.000000 Done]
[DOWNLOAD_STATUS pi-billion.txt 91273edf-e095-4ff4-8c3c-128523698bd6 : 99.000000 Done]
[DOWNLOAD_STATUS pi-billion.txt 91273edf-e095-4ff4-8c3c-128523698bd6 : 100.000000 Done]

DOWNLOAD_COMPLETE pi-billion.txt 91273edf-e095-4ff4-8c3c-128523698bd6]
Download complete: pi-billion.txt

```
`This would run goget with 100 MB reserved for downloading objects`

TODO:

  1. Download Directories and sub directories
  2. Add a Tracker for downloads
