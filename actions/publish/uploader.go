package main

import (
	"github.com/prometheus/common/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

func main() {
	// get version from env
	version := "0.0.0"
	s3Folder := "path/to/s3"
	artifactsStore := "path/to/downloads"
	fileMappingConfig := "path/to/config"
	binaryName := "some-app-name"

	// take schema
	fileContent, err := readFileContent("filepath")

	if err != nil {
		log.Fatal(err)
	}

	//handle err
	// parse schema
	config := parseConfig(fileContent)
	// iterate over config
	uploadArtifacts(binaryName, config, version)
	// --> replace variables
	// --> copy (move?) file
}

type uploadArtifactConfig struct {
	src  string   `json:"src"`
	dest string   `json:"dest"`
	arch []string `json:"arch"`
}

type uploadConfig []uploadArtifactConfig

func readFileContent(filePath string) ([]byte, error) {
	fileContent, err := ioutil.ReadFile(filePath)

	return fileContent, err
}

func parseConfig(fileContent []byte) (uploadConfig, error) {

	var config uploadConfig

	err := yaml.Unmarshal(fileContent, &config)

	if err != nil {
		return uploadConfig{}, err

	}

	return config, nil
}

func uploadArtifact(binaryName string, config uploadArtifactConfig, version string) {

	if len(config.arch) > 0 {

		for _, arch := range config.arch {
			srcPath, destPath := replacePlaceholders(config.src, config.dest, arch, binaryName, version)

			//srcPath = os.Jo "" + srcPath

			input, err := ioutil.ReadFile(srcPath)
			if err != nil {
				log.Fatal(err)
				return
			}

			err = ioutil.WriteFile(destPath, input, 0644)
			if err != nil {
				log.Error("Error creating", destPath)
				log.Fatal(err)
				return
			}

		}
	} else {
		replacePlaceholders(config.src, config.dest, "", binaryName, version)
	}
}

func uploadArtifacts(binaryName string, config uploadConfig, version string) {
	for _, artifactConfig := range config {
		uploadArtifact(binaryName, artifactConfig, version)
	}
}

func replacePlaceholders(srcTemplate, destTemplate, arch, binaryName, version string) (string, string) {

	srcFileName := ""
	destPath := ""
	return srcFileName, destPath
}
