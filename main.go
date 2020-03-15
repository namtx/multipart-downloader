package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/cheggaaa/pb"
)

var wg sync.WaitGroup

const (
	AcceptRanges     = "Accept-Ranges"
	ContentLength    = "Content-Length"
	Range            = "Range"
	partLength       = 2 * 1000 * 1000
	numberGoroutines = 8
)

func main() {
	wg.Add(numberGoroutines)

	url := os.Args[1]
	fileName := GetFileName(url)

	// send HEAD request
	response, err := http.Head(url)
	checkError(err)

	if response.Header[AcceptRanges][0] != "bytes" {
		fmt.Printf("%s does NOT support Multipart Download\n", url)
		os.Exit(1)
	}

	contentLength, err := strconv.ParseInt(response.Header.Get(ContentLength), 10, 64)
	checkError(err)

	fmt.Printf("Total: %d\n", contentLength)

	parts := CalculateNumberOfParts(contentLength, partLength)
	tasks := make(chan Task, parts)
	client := &http.Client{}
	progressBars := make([]*pb.ProgressBar, 0)

	for i := 0; i < numberGoroutines; i++ {
		go Download(tasks, client, url, fileName)
	}

	files := make([]string, 0)

	for part := 0; part < parts; part++ {
		partFileName := fmt.Sprintf("%s-%d", fileName, part)
		progressBar := pb.New(partLength).Prefix(partFileName)

		task := Task{
			Part:        part,
			ProgressBar: progressBar,
		}

		tasks <- task

		progressBars = append(progressBars, progressBar)
		files = append(files, partFileName)
	}

	progressBarPool, err := pb.StartPool(progressBars...)
	checkError(err)

	close(tasks)

	wg.Wait()

	progressBarPool.Stop()

	err = Join(files, fileName)

	checkError(err)
}
