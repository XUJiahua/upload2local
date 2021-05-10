package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	clientUsage = `client usage: curl -F "data=@main.go" localhost:8090`
	outputDir   = "/tmp"
)

func receive(w http.ResponseWriter, r *http.Request) {
	var (
		mediaType string
		params    map[string]string
		err       error
	)

	defer func() {
		if err != nil {
			fmt.Fprintf(w, "err: %v\n", err)
			fmt.Fprintf(w, "%s\n", clientUsage)
		} else {
			fmt.Fprintf(w, "ok\n")
		}
	}()

	mediaType, params, err = mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		return
	}
	if !strings.HasPrefix(mediaType, "multipart/") {
		err = errors.New("media type is not prefixed with multipart/")
		return
	}
	mr := multipart.NewReader(r.Body, params["boundary"])
	for {
		var (
			p       *multipart.Part
			content []byte
		)
		p, err = mr.NextPart()
		if err == io.EOF {
			err = nil
			break
		}
		if err != nil {
			return
		}
		content, err = io.ReadAll(p)
		if err != nil {
			return
		}
		err = writeFile(p.FileName(), content)
		if err != nil {
			return
		}
	}
}

func writeFile(filename string, content []byte) error {
	filename = filepath.Join(outputDir, filename)
	return ioutil.WriteFile(filename, content, 0666)
}

func main() {
	http.HandleFunc("/", receive)

	handleErr(http.ListenAndServe(":8090", nil))
}

func handleErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}
