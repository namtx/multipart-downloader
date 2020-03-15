package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func Download(tasks chan Task, client *http.Client, url string, fileName string) {
	defer wg.Done()

	for {
		task, ok := <-tasks

		if !ok {
			return
		}

		request, err := http.NewRequest("GET", url, nil)
		checkError(err)

		rangeHeader := fmt.Sprintf("bytes=%d-%d", task.Part*partLength, (task.Part+1)*partLength-1)
		request.Header.Add(Range, rangeHeader)

		response, err := client.Do(request)
		checkError(err)

		outFileName := fmt.Sprintf("%s-%d", fileName, task.Part)
		filePath := filepath.Join("/tmp", filepath.Base(outFileName))

		outFile, err := os.Create(filePath)

		defer outFile.Close()

		// Progress bar
		barReader := task.ProgressBar.NewProxyReader(response.Body)
		defer barReader.Close()
		io.Copy(outFile, barReader)

		task.ProgressBar.Finish()
	}
}
