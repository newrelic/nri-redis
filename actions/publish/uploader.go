package main

import (
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"path"
	"strings"
)

const (
	placeholderForVersion = "{version}"
	placeholderForArch    = "{arch}"
	placeholderForAppName = "{app_name}"
)

type config struct {
	version              string
	artifactsDestFolder  string
	artifactsSrcFolder   string
	uploadSchemaFilePath string
	appName              string
}

type uploadArtifactSchema struct {
	Src  string   `yaml:"src"`
	Dest string   `yaml:"dest"`
	Arch []string `yaml:"arch"`
}

type uploadArtifactsSchema []uploadArtifactSchema

func main() {

	viper.BindEnv("version")
	viper.BindEnv("artifactsDestFolder")
	viper.BindEnv("artifactsSrcFolder")
	viper.BindEnv("uploadSchemaFilePath")
	viper.BindEnv("appName")

	pflag.String("version", "0.0.0", "asset version")
	pflag.String("artifactsDestFolder", "", "artifacts destination folder")
	pflag.String("artifactsSrcFolder", "", "artifacts source folder")
	pflag.String("uploadSchemaFilePath", "", "artifacts source folder")
	pflag.String("appName", "", "app name")

	pflag.Parse()

	viper.BindPFlags(pflag.CommandLine)

	fmt.Println(viper.GetString("version"))

	conf := config{
		version:              viper.GetString("version"),
		artifactsDestFolder:  viper.GetString("artifactsDestFolder"),
		artifactsSrcFolder:   viper.GetString("artifactsSrcFolder"),
		uploadSchemaFilePath: viper.GetString("uploadSchemaFilePath"),
		appName:              viper.GetString("appName"),
	}

	uploadSchemaContent, err := readFileContent(conf.uploadSchemaFilePath)

	if err != nil {
		log.Fatal(err)
	}

	uploadSchema, err := parseUploadSchema(uploadSchemaContent)

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

	if len(schema.Arch) > 0 {

		for _, arch := range schema.Arch {
			srcPath, destPath := replacePlaceholders(schema.Src, schema.Dest, arch, conf.appName, conf.version)

			srcPath = path.Join(conf.artifactsSrcFolder, srcPath)
			destPath = path.Join(conf.artifactsDestFolder, destPath)

			input, err := ioutil.ReadFile(srcPath)
			if err != nil {
				log.Fatal(err)
				return
			}

			err = ioutil.WriteFile(destPath, input, 0644)
			if err != nil {
				log.Print("Error creating", destPath)
				log.Fatal(err)
				return
			}

		}
	} else {
		replacePlaceholders(schema.Src, schema.Dest, "", conf.appName, conf.version)
	}
}

func uploadArtifacts(conf config, schema uploadArtifactsSchema) {
	for _, artifactSchema := range schema {
		uploadArtifact(conf, artifactSchema)
	}
}

func replacePlaceholders(srcTemplate, destTemplate, arch, binaryName, version string) (string, string) {

	srcFileName := strings.Replace(srcTemplate, placeholderForVersion, version, 1)
	srcFileName = strings.Replace(srcFileName, placeholderForArch, arch, 1)
	srcFileName = strings.Replace(srcFileName, placeholderForAppName, binaryName, 1)

	destPath := strings.Replace(destTemplate, placeholderForVersion, version, 1)
	destPath = strings.Replace(destPath, placeholderForArch, arch, 1)
	destPath = strings.Replace(destPath, placeholderForAppName, binaryName, 1)

	return srcFileName, destPath
}
