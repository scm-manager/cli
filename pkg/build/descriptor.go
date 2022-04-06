package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type MetadataJson struct {
	Tag  string
	Date string
}

type Artifact struct {
	Name   string
	Type   string
	Goarch string
	Goos   string
}

type Package struct {
	Type     string
	Os       string `yaml:",omitempty"`
	Arch     string `yaml:",omitempty"`
	Url      string `yaml:",omitempty"`
	Checksum string `yaml:",omitempty"`
}

type Descriptor struct {
	Tag      string
	Date     string
	Packages []Package
}

func main() {
	dir := os.Args[1]
	file, err := ioutil.ReadFile(path.Join(dir, "metadata.json"))
	if err != nil {
		log.Fatal("could not read metadata.json")
	}

	metadata := MetadataJson{}
	err = json.Unmarshal(file, &metadata)
	if err != nil {
		log.Fatal("could not unmarshal metadata.json")
	}

	file, err = ioutil.ReadFile(path.Join(dir, "artifacts.json"))
	if err != nil {
		log.Fatal("could not read artifacts.json")
	}

	var artifacts []Artifact
	err = json.Unmarshal(file, &artifacts)
	if err != nil {
		log.Fatal("could not unmarshal artifacts.json")
	}

	var packages []Package
	checksums := readChecksums(dir)
	for _, a := range artifacts {
		if a.Type == "Archive" {
			checksum := checksums[a.Name]
			url := fmt.Sprintf("https://packages.scm-manager.org/repository/scm-cli-releases/%s/%s", metadata.Tag, a.Name)
			packages = append(packages, Package{
				Os:       a.Goos,
				Arch:     a.Goarch,
				Checksum: checksum,
				Type:     filepath.Ext(a.Name)[1:],
				Url:      url,
			})
		} else if a.Type == "Linux Package" {
			packages = append(packages, Package{Os: a.Goos, Arch: a.Goarch, Type: filepath.Ext(a.Name)[1:]})
		} else if a.Type == "Homebrew" {
			packages = append(packages, Package{Type: "homebrew"})
		} else if a.Type == "Scoop Manifest" {
			packages = append(packages, Package{Type: "scoop"})
		}
	}

	date, err := time.Parse(time.RFC3339Nano, metadata.Date)
	if err != nil {
		log.Fatal("could not parse release date")
	}
	descriptor := Descriptor{Tag: metadata.Tag, Date: date.UTC().Format(time.RFC3339), Packages: packages}

	bytes, err := yaml.Marshal(&descriptor)
	if err != nil {
		log.Fatal("could not unmarshal artifacts.json")
	}
	fmt.Println(string(bytes))

}

func readChecksums(dir string) map[string]string {
	checksumFile, err := os.Open(path.Join(dir, "checksums.txt"))
	if err != nil {
		log.Fatal("could not read checksums.txt")
	}
	defer checksumFile.Close()
	scanner := bufio.NewScanner(checksumFile)
	checksums := make(map[string]string)
	for scanner.Scan() {
		line := scanner.Text()
		result := strings.Split(line, "  ")
		checksums[result[1]] = result[0]
	}
	err = scanner.Err()
	if err != nil {
		log.Fatal("error during scanning checksums.txt")
	}
	return checksums
}
