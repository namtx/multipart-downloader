package main

import (
	"fmt"
	"math"
	"os"
	"strings"
)

func GetFileName(url string) string {
	splittedURL := strings.Split(url, "/")

	return splittedURL[len(splittedURL)-1]
}

func CalculateNumberOfParts(contentLength int64, partLength int) int {
	return int(math.Ceil(float64(contentLength) / float64(partLength)))
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err.Error())
		os.Exit(1)
	}
}
