package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type NameMap struct {
	FileName string `json:"file"`
}

type ScmmUploadRequest struct {
	Message  string  `json:"commitMessage"`
	Branch   string  `json:"branch"`
	Revision string  `json:"expectedRevision,omitempty"`
	Names    NameMap `json:"names,omitempty"`
}

func main() {
	username := os.Getenv("SCMM_USERNAME")
	if username == "" {
		log.Fatal("missing SCMM username")
	}
	password := os.Getenv("SCMM_PASSWORD")
	if password == "" {
		log.Fatal("missing SCMM password")
	}

	args := os.Args

	source := args[1]
	repo := args[2]
	branch := args[3]
	path := args[4]
	commitMessage := args[5]

	url := fmt.Sprintf("https://ecosystem.cloudogu.com/scm/api/v2/edit/scm-manager/%s/create/%s", repo, filepath.Dir(path))

	uploadFile := mustOpen(source)

	uploadRequest := ScmmUploadRequest{
		Message: commitMessage,
		Branch:  branch,
		Names:   NameMap{FileName: filepath.Base(path)},
	}

	commit, err := json.Marshal(&uploadRequest)
	if err != nil {
		log.Fatal("Could not marshal json")
	}

	values := map[string]io.Reader{
		"commit": bytes.NewReader(commit),
		"file":   uploadFile,
	}
	err = Upload(url, values, username, password)
	if err != nil {
		log.Fatal("could not send upload request", err)
	}
}

func Upload(url string, values map[string]io.Reader, username string, password string) (err error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range values {
		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		// Add an image file
		if _, ok := r.(*os.File); ok {
			if fw, err = w.CreateFormFile(key, "file"); err != nil {
				return
			}
		} else {
			// Add other fields
			if fw, err = w.CreateFormField(key); err != nil {
				return
			}
		}
		if _, err = io.Copy(fw, r); err != nil {
			return err
		}

	}
	w.Close()

	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Authorization", "Basic "+basicAuth(username, password))

	// Submit the request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	if res.StatusCode != http.StatusCreated {
		buf := new(strings.Builder)
		_, _ = io.Copy(buf, res.Body)
		log.Print(buf.String())
		err = fmt.Errorf("bad status: %s", res.Status)
	}
	return
}

func mustOpen(f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		panic(err)
	}
	return r
}

func basicAuth(username string, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
