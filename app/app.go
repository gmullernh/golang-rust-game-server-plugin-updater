package main

import (
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const Url = `https://umod.org/plugins/search.json?query={0}&page=1&sort=title&sortdir=asc`

func main() {

	workingDir, err := os.Getwd()
	if err != nil {
		fmt.Println(`Error getting working directory:`, err)
		return
	}
	
	// Create the relative path to the 'data' folder
	folderPath := filepath.Join(workingDir, `/plugins`)

	files, err := os.ReadDir(folderPath)
	if err != nil {
		fmt.Println(`Error reading directory:`, err)
		return
	}

	for _, file := range files {
		// Check if it's a regular file (not a directory)
		if !file.IsDir() {
			filename := file.Name()

			myUrl := getDownloadLink(filename)

			if myUrl != `` {
				err := downloadFile(myUrl, folderPath+`/`+filename)
				if err != nil {
					fmt.Println(`error`, err)
					return
				}
				fmt.Println(`found`, filename, `@`, myUrl)
			} else {
				fmt.Println(`not found`, filename)
			}

			waitForSeconds(5)
		}
	}
}

func waitForSeconds(seconds time.Duration) {
	duration := seconds * time.Second
	time.Sleep(duration)
}

func downloadFile(urlPath, filepath string) error {

	if len(urlPath) == 0 {
		return nil
	}

	fixUrl := strings.Replace(urlPath, `/`, ``, -1)
	fixUrl = strings.Replace(fixUrl, `\`, `/`, -1)

	response, err := http.Get(fixUrl)
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(`error while downloading json data`)
		}
	}(response.Body)

	// Create the local file
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(`couldn't create the file`)
		}
	}(file)

	// Copy the content from the response body to the local file
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}

func getDownloadLink(filename string) string {

	escapedUrl := url.QueryEscape(filename)
	escapedUrl = escapedUrl[:len(escapedUrl)-3]
	myUrl := strings.Replace(Url, `{0}`, escapedUrl, -1)

	response, err := http.Get(myUrl)
	if err != nil {
		fmt.Println(`Error:`, err)
		return ``
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(`Error while fetching content.`)
		}
	}(response.Body)

	// Read the response body as a string
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(`Error reading response:`, err)
		return ``
	}
	data := string(bodyBytes)

	// get all download urls from JSON response
	download := gjson.Get(data, `data.#.download_url`)

	// Check if the filename is equal to the download url filename
	for i := 0; i < len(download.Array()); i++ {
		dnwUrl := strings.Replace((download.Array()[i]).Raw, `"`, ``, -1)

		if strings.Contains(dnwUrl, `/`+filename) {
			return dnwUrl
		}
	}

	return ``
}
