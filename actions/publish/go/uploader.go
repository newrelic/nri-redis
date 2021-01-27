package main

import (
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
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
	viper.BindEnv("artifacts_dest_folder")
	viper.BindEnv("artifacts_src_folder")
	viper.BindEnv("uploadSchema_file_path")
	viper.BindEnv("app_name")

	pflag.String("version", "0.0.0", "asset version")
	pflag.String("artifactsDestFolder", "", "artifacts destination folder")
	pflag.String("artifactsSrcFolder", "", "artifacts source folder")
	pflag.String("uploadSchemaFilePath", "", "upload schema file path")
	pflag.String("appName", "", "app name")

	pflag.Parse()

	viper.BindPFlags(pflag.CommandLine)

	getFirstNotEmpty := func(first, second string) string {
		if first != "" {
			return first
		}

		return second
	}

	conf := config{
		version:              viper.GetString("version"),
		artifactsDestFolder:  getFirstNotEmpty(viper.GetString("artifactsDestFolder"), viper.GetString("artifacts_dest_folder")),
		artifactsSrcFolder:   getFirstNotEmpty(viper.GetString("artifactsSrcFolder"), viper.GetString("artifacts_src_folder")),
		uploadSchemaFilePath: getFirstNotEmpty(viper.GetString("uploadSchemaFilePath"), viper.GetString("uploadSchema_file_path")),
		appName:              getFirstNotEmpty(viper.GetString("appName"), viper.GetString("app_name")),
	}

	log.Println(fmt.Sprintf("%v", conf))

	uploadSchemaContent, err := readFileContent(conf.uploadSchemaFilePath)

	if err != nil {
		log.Fatal(err)
	}

	uploadSchema, err := parseUploadSchema(uploadSchemaContent)

	if err != nil {
		log.Fatal(err)
	}

	err = uploadArtifacts(conf, uploadSchema)

	if err != nil {
		log.Fatal(err)
	}
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

func uploadArtifact(conf config, schema uploadArtifactSchema) error {

	if len(schema.Arch) > 0 {

		for _, arch := range schema.Arch {
			srcPath, destPath := replacePlaceholders(schema.Src, schema.Dest, arch, conf.appName, conf.version)

			srcPath = path.Join(conf.artifactsSrcFolder, srcPath)
			destPath = path.Join(conf.artifactsDestFolder, destPath)

			destDirectory := filepath.Dir(destPath)

			if _, err := os.Stat(destDirectory); os.IsNotExist(err) {
				// set right permissions
				os.Mkdir(destDirectory, 0644)
			}

			input, err := ioutil.ReadFile(srcPath)
			if err != nil {
				return err
			}

			err = ioutil.WriteFile(destPath, input, 0644)
			if err != nil {
				return err
			}

		}
	} else {
		replacePlaceholders(schema.Src, schema.Dest, "", conf.appName, conf.version)
	}
	return nil
}

func uploadArtifacts(conf config, schema uploadArtifactsSchema) error {
	for _, artifactSchema := range schema {
		err := uploadArtifact(conf, artifactSchema)
		if err != nil{
			return err
		}
	}
	return nil
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
