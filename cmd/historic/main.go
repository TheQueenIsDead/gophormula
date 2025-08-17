package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gophormula/pkg/livetiming/messages"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
)

const (
	BOM = "\xEF\xBB\xBF"
)

func GetLiveTimingFile[T any](url *url.URL, out *T) error {

	// HTTP GET the file
	res, err := http.Get(url.String())
	if err != nil || res.StatusCode != 200 {
		fmt.Println("unable to get", url.String())
		return err
	}

	// Read out the body and remove any UTF8 BOM's
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("unable to read body")
		return err
	}
	body = bytes.TrimPrefix(body, []byte(BOM))

	// Marshall the response into the given struct
	err = json.Unmarshal(body, out)
	if err != nil {
		fmt.Println("unable to unmarshal")
		return err
	}

	return nil
}

func main() {

	if len(os.Args) < 3 {
		fmt.Println("Usage: historic <url> <out>")
	}

	in, out := os.Args[1], os.Args[2]

	base, err := url.Parse(in)
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll(strings.ReplaceAll(out+base.Path, "static/", ""), 0750)
	if err != nil {
		panic(err)
	}

	fmt.Println("Retrieving index...")
	var index messages.Index
	err = GetLiveTimingFile(base.JoinPath("Index.json"), &index)
	if err != nil {
		panic(err)
	}

	var files []*url.URL
	for _, feed := range index.GetFeeds() {
		if feed.KeyFramePath != "" {
			files = append(files, base.JoinPath(feed.KeyFramePath))
		}
		if feed.StreamPath != "" {
			files = append(files, base.JoinPath(feed.StreamPath))
		}
	}

	fmt.Println("Retrieving race files...")
	var wg sync.WaitGroup
	for _, file := range files {
		wg.Add(1)
		file := file
		go func() {
			defer wg.Done()

			// Get
			res, err := http.Get(file.String())
			if err != nil || res.StatusCode != 200 {
				fmt.Println(err, file)
			}

			// Read
			body, err := io.ReadAll(res.Body)
			if err != nil {
				fmt.Println("unable to read body")
			}

			// Write
			path := strings.ReplaceAll(out+res.Request.URL.Path, "static/", "")
			f, err := os.Create(path)
			if err != nil {
				fmt.Println("unable to create file:", path, file)
			}
			defer f.Close()
			_, err = f.Write(body)
			if err != nil {
				fmt.Println("unable to write file")
			}
		}()
	}
	wg.Wait()
	fmt.Println("Race files saved.")
}
