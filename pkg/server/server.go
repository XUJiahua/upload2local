package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/xujiahua/upload2local/pkg/model"
	"github.com/xujiahua/upload2local/pkg/split"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Server struct {
	// for saving merged file
	workDir string
	// for saving file parts
	tmpDir string
	port   int
}

func New(workDir string, port int) (*Server, error) {
	err := os.MkdirAll(workDir, 0777)
	if err != nil {
		return nil, err
	}

	return &Server{
		workDir: workDir,
		tmpDir:  os.TempDir(),
		port:    port,
	}, nil
}

func (s Server) receive(w http.ResponseWriter, r *http.Request) {
	var (
		mediaType string
		params    map[string]string
		err       error
		res       model.Response
	)

	defer func() {
		if err != nil {
			logrus.Error(err)
			res.Code = model.CodeFailure
			res.Message = err.Error()
		} else {
			res.Code = model.CodeSuccess
		}

		err := json.NewEncoder(w).Encode(res)
		if err != nil {
			logrus.Error(err)
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

		filename := rebase(p.FileName(), s.tmpDir)
		err = ioutil.WriteFile(filename, content, 0666)
		if err != nil {
			return
		}
		logrus.Infof("received part file %s, with %d bytes\n", filename, len(content))
	}
}

func (s Server) complete(w http.ResponseWriter, r *http.Request) {
	var err error
	var res model.Response
	defer func() {
		if err != nil {
			logrus.Error(err)
			res.Code = model.CodeFailure
			res.Message = err.Error()
		} else {
			res.Code = model.CodeSuccess
		}

		err := json.NewEncoder(w).Encode(res)
		if err != nil {
			logrus.Error(err)
		}
	}()

	var complete model.CompleteRequest
	err = json.NewDecoder(r.Body).Decode(&complete)
	if err != nil {
		return
	}
	if complete.Filename == "" || len(complete.PartFiles) == 0 {
		err = errors.New("request is invalid")
		return
	}

	for i := 0; i < len(complete.PartFiles); i++ {
		complete.PartFiles[i] = rebase(complete.PartFiles[i], s.tmpDir)
	}
	complete.Filename = rebase(complete.Filename, s.workDir)
	err = split.Merge(complete.PartFiles, complete.Filename)
}

func rebase(filename, newDir string) string {
	return filepath.Join(newDir, filepath.Base(filename))
}

func (s Server) Start() error {
	http.HandleFunc("/", s.receive)
	http.HandleFunc("/complete", s.complete)

	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil)
}
