package main

import (
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const (
	placeholderForVersion = "{version}"
	placeholderForArch    = "{arch}"
	placeholderForAppName = "{app_name}"
	placeholderForSrc     = "{src}"
)

type config struct {
	version              string
	tag                  string
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

	viper.BindEnv("tag")
	viper.BindEnv("artifacts_dest_folder")
	viper.BindEnv("artifacts_src_folder")
	viper.BindEnv("uploadSchema_file_path")
	viper.BindEnv("app_name")

	pflag.String("tag", "v0.0.0", "asset git tag")
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
		version:              strings.Replace(viper.GetString("tag"), "v", "", -1),
		tag:                  viper.GetString("tag"),
		artifactsDestFolder:  getFirstNotEmpty(viper.GetString("artifactsDestFolder"), viper.GetString("artifacts_dest_folder")),
		artifactsSrcFolder:   getFirstNotEmpty(viper.GetString("artifactsSrcFolder"), viper.GetString("artifacts_src_folder")),
		uploadSchemaFilePath: getFirstNotEmpty(viper.GetString("uploadSchemaFilePath"), viper.GetString("uploadSchema_file_path")),
		appName:              getFirstNotEmpty(viper.GetString("appName"), viper.GetString("app_name")),
	}

	log.Println(fmt.Sprintf("config: %v", conf))

	uploadSchemaContent, err := readFileContent(conf.uploadSchemaFilePath)

	if err != nil {
		log.Fatal(err)
	}

	uploadSchema, err := parseUploadSchema(uploadSchemaContent)

	if err != nil {
		log.Fatal(err)
	}

	err = downloadArtifacts(conf, uploadSchema)

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

func downloadArtifact(conf config, schema uploadArtifactSchema) error {

	if len(schema.Arch) > 0 {

		for _, arch := range schema.Arch {
			downloadUrl := ""

			srcPath, _ := replacePlaceholders(schema.Src, schema.Dest, arch, conf.appName, conf.version)

			srcPath = path.Join(conf.artifactsSrcFolder, srcPath)

			log.Println("[ ] Download " + downloadUrl + " into " + srcPath)

			// add download here

			log.Println("[✔] Download " + downloadUrl + " into " + srcPath)

		}
	} else {
		replacePlaceholders(schema.Src, schema.Dest, "", conf.appName, conf.version)
	}
	return nil
}

func downloadArtifacts(conf config, schema uploadArtifactsSchema) error {
	for _, artifactSchema := range schema {
		err := downloadArtifact(conf, artifactSchema)
		if err != nil {
			return err
		}
	}
	return nil
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
				err = os.MkdirAll(destDirectory, 0777)
				if err != nil {
					return err
				}
			}

			log.Println("[ ] Copy " + srcPath + " into " + destPath)

			input, err := ioutil.ReadFile(srcPath)
			if err != nil {
				return err
			}

			err = ioutil.WriteFile(destPath, input, 0777)
			if err != nil {
				return err
			}

			log.Println("[✔] Copy " + srcPath + " into " + destPath)

		}
	} else {
		replacePlaceholders(schema.Src, schema.Dest, "", conf.appName, conf.version)
	}
	return nil
}

func uploadArtifacts(conf config, schema uploadArtifactsSchema) error {
	for _, artifactSchema := range schema {
		err := uploadArtifact(conf, artifactSchema)
		if err != nil {
			return err
		}
	}
	return nil
}

func replacePlaceholders(srcTemplate, destTemplate, arch, appName, version string) (string, string) {

	srcFileName := strings.Replace(srcTemplate, placeholderForVersion, version, -1)
	srcFileName = strings.Replace(srcFileName, placeholderForArch, arch, -1)
	srcFileName = strings.Replace(srcFileName, placeholderForAppName, appName, -1)

	destPath := strings.Replace(destTemplate, placeholderForVersion, version, -1)
	destPath = strings.Replace(destPath, placeholderForArch, arch, -1)
	destPath = strings.Replace(destPath, placeholderForAppName, appName, -1)
	destPath = strings.Replace(destPath, placeholderForSrc, srcFileName, -1)

	return srcFileName, destPath
}

// TAG = v1.1.1
// VERSION = 1.1.1 (TAG.repalce(v,""))

func generateDownloadUrl(srcTemplate, arch, binaryName, version string) {
	//template := "https://github.com/{app_name}/releases/download/${GH_TAG}/${PKG_NAME}"
}

func downloadFile(url, fileName string) error {

	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return fmt.Errorf("error on download: %v", err)
	}

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}
