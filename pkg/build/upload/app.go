package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Person struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type GitHubUploadRequest struct {
	Message string `json:"message"`
	Branch  string `json:"branch"`
	Author  Person `json:"author"`
	Content string `json:"content"`
}

func main() {
	token := os.Getenv("GITHUB_API_TOKEN")
	if token == "" {
		log.Fatal("missing GitHub API Token")
	}

	args := os.Args

	source := args[1]
	repo := args[2]
	branch := args[3]
	path := args[4]
	commitMessage := args[5]

	url := fmt.Sprintf("https://api.github.com/repos/scm-manager/%s/contents/%s", repo, path)

	file, err := ioutil.ReadFile(source)
	if err != nil {
		log.Fatal("could not read source file")
	}
	content := base64.StdEncoding.EncodeToString(file)

	uploadRequest := GitHubUploadRequest{
		Message: commitMessage,
		Branch:  branch,
		Author:  Person{Name: "CES Marvin", Email: "cesmarvin@cloudogu.com"},
		Content: content,
	}
	data, err := json.Marshal(&uploadRequest)
	if err != nil {
		log.Fatal("Could not marshal json")
	}
	request, err := http.NewRequest("PUT", url, bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("could not create upload request")
	}
	request.Header.Set("Authorization", "Token "+token)
	request.Header.Set("Accept", "application/vnd.github.v3+json")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal("could not send upload request")
	}

	if response.StatusCode >= 300 {
		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal("Could not read response body")
		}

		log.Fatalf("upload request failed: %d\n\n%s", response.StatusCode, string(body))
	}

}
