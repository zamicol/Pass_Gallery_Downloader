package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

var (
	count   = 0
	failed  = 0
	success = 0
	// Start is where the counting should start, stop is where to stop the program.
	start = 679224233
	stop  = 679324233
	total = stop - start

	//URL base will be in this format: https://example.com/pictures/13/163/13163938/lowres/
	urlBase = "https://example.com/pictures/13/163/13163938/lowres/"
)

func main() {
	// Bytes that begin a page that doesn't get a jpg hit
	var badPage = []byte{239, 187, 191, 60, 63, 120, 109, 108, 32, 118, 101, 114, 115, 105}

	for i := start; i <= stop; i-- {
		count++
		var filename string = strconv.Itoa(i) + ".jpg"
		fmt.Println("Attepting: " + filename)

		fileUrl := urlBase + filename

		resp, err := http.Get(fileUrl)
		if err != nil {
			fmt.Println("Error opeing page.")
		}
		defer resp.Body.Close()
		// reads html as a slice of bytes
		html, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("There was an error!")
		}
		if bytes.HasPrefix(html, badPage) {
			failed += 1
			continue
		}

		fmt.Println("Page Hit!")
		if err := DownloadFile("photos/"+filename, fileUrl); err != nil {
			fmt.Println("failed")
			failed += 1
			continue
		}
		success += 1
		fmt.Println("Success:  " + filename)
		stats()
	}
	fmt.Println("Finished downloading everything.")
	stats()
}

func stats() {
	fmt.Printf("Total Scaned: %v\n", count)
	fmt.Printf("Total Failed: %v\n", failed)
	fmt.Printf("Total Success: %v\n", success)
	fmt.Printf("Total left to scan: %v \n", total-count)
}

// DownloadFile downloads the a file to disk from url.
func DownloadFile(filepath string, url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
