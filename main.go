package main

import "fmt"

func main() {
	url := "https://raw.githubusercontent.com/unsubd/geektrust-family/master/input.txt"
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
