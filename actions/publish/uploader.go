package main

import (
	"github.com/prometheus/common/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type config struct {
	version               string
	artifactsDestFolder   string
	artifactsSourceFolder string
	uploadSchemaFilePath  string
	binaryName            string
}

type uploadArtifactSchema struct {
	src  string   `json:"src"`
	dest string   `json:"dest"`
	arch []string `json:"arch"`
}

type uploadArtifactsSchema []uploadArtifactSchema

func main() {

	conf := config{
		version:               "0.0.0",
		artifactsDestFolder:   "path/to/s3",
		artifactsSourceFolder: "path/to/downloads",
		uploadSchemaFilePath:  "path/to/config",
		binaryName:            "some-app-name",
	}

	fileContent, err := readFileContent("filepath")

	if err != nil {
		log.Fatal(err)
	}

	uploadSchema, err := parseUploadSchema(fileContent)

	if err != nil {
		log.Fatal(err)
	}

	uploadArtifacts(conf, uploadSchema)
}

func readFileContent(filePath string) ([]byte, error) {
	fileContent, err := ioutil.ReadFile(filePath)

	return fileContent, err
}

func parseUploadSchema(fileContent []byte) (uploadArtifactsSchema, error) {

	var schema uploadArtifactsSchema

	err := yaml.Unmarshal(fileContent, &schema)

	if err != nil {
		return uploadArtifactsSchema{}, err

	}

	return schema, nil
}

func uploadArtifact(conf config, schema uploadArtifactSchema) {

	if len(schema.arch) > 0 {

		for _, arch := range schema.arch {
			srcPath, destPath := replacePlaceholders(schema.src, schema.dest, arch, conf.binaryName, conf.version)

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
		replacePlaceholders(schema.src, schema.dest, "", conf.binaryName, conf.version)
	}
}

func uploadArtifacts(conf config, schema uploadArtifactsSchema) {
	for _, artifactSchema := range schema {
		uploadArtifact(conf, artifactSchema)
	}
}

func replacePlaceholders(srcTemplate, destTemplate, arch, binaryName, version string) (string, string) {

	srcFileName := ""
	destPath := ""
	return srcFileName, destPath
}
