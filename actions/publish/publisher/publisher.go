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
	placeholderForRepoName = "{repo_name}"
	placeholderForAppName  = "{app_name}"
	placeholderForArch     = "{arch}"
	placeholderForTag      = "{tag}"
	placeholderForVersion  = "{version}"
	placeholderForSrc      = "{src}"
	urlTemplate            = "https://github.com/{repo_name}/releases/download/{tag}/{src}"
)

type config struct {
	repoName             string
	appName              string
	tag                  string
	version              string
	artifactsDestFolder  string
	artifactsSrcFolder   string
	uploadSchemaFilePath string
}

type uploadArtifactSchema struct {
	Src  string   `yaml:"src"`
	Dest string   `yaml:"dest"`
	Arch []string `yaml:"arch"`
}

type uploadArtifactsSchema []uploadArtifactSchema

func main() {
	conf := loadConfig()
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

func loadConfig() config {
	// TODO: make all the config required
	viper.BindEnv("repo_name")
	viper.BindEnv("app_name")
	viper.BindEnv("tag")
	viper.BindEnv("artifacts_dest_folder")
	viper.BindEnv("artifacts_src_folder")
	viper.BindEnv("uploadSchema_file_path")

	pflag.String("repoName", "", "repo name")
	pflag.String("appName", "", "app name")
	pflag.String("tag", "", "asset git tag")
	pflag.String("artifactsDestFolder", "", "artifacts destination folder")
	pflag.String("artifactsSrcFolder", "", "artifacts source folder")
	pflag.String("uploadSchemaFilePath", "", "upload schema file path")

	pflag.Parse()

	viper.BindPFlags(pflag.CommandLine)

	getFirstNotEmpty := func(first, second string) string {
		if first != "" {
			return first
		}

		return second
	}

	return config{
		repoName:             getFirstNotEmpty(viper.GetString("repoName"), viper.GetString("repo_name")),
		appName:              getFirstNotEmpty(viper.GetString("appName"), viper.GetString("app_name")),
		tag:                  viper.GetString("tag"),
		version:              strings.Replace(viper.GetString("tag"), "v", "", -1),
		artifactsDestFolder:  getFirstNotEmpty(viper.GetString("artifactsDestFolder"), viper.GetString("artifacts_dest_folder")),
		artifactsSrcFolder:   getFirstNotEmpty(viper.GetString("artifactsSrcFolder"), viper.GetString("artifacts_src_folder")),
		uploadSchemaFilePath: getFirstNotEmpty(viper.GetString("uploadSchemaFilePath"), viper.GetString("uploadSchema_file_path")),
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

	for i, _ := range schema {
		if schema[i].Arch == nil {
			schema[i].Arch = []string{""}
		}
	}

	return schema, nil
}

func downloadArtifact(conf config, schema uploadArtifactSchema) error {
	for _, arch := range schema.Arch {
		srcFile, _ := replaceSrcDestTemplates(schema.Src, schema.Dest, conf.repoName, conf.appName, arch, conf.tag, conf.version)
		url := generateDownloadUrl(urlTemplate, conf.repoName, conf.appName, arch, conf.tag, conf.version, srcFile)

		destPath := path.Join(conf.artifactsSrcFolder, srcFile)

		log.Println("[ ] Download " + url + " into " + destPath)

		err := downloadFile(url, destPath)
		if err != nil {
			return err
		}

		log.Println("[✔] Download " + url + " into " + destPath)
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
	for _, arch := range schema.Arch {
		srcPath, destPath := replaceSrcDestTemplates(schema.Src, schema.Dest, conf.repoName, conf.appName, arch, conf.tag, conf.version)

		srcPath = path.Join(conf.artifactsSrcFolder, srcPath)
		destPath = path.Join(conf.artifactsDestFolder, destPath)

		destDirectory := filepath.Dir(destPath)

		if _, err := os.Stat(destDirectory); os.IsNotExist(err) {
			// set right permissions
			err = os.MkdirAll(destDirectory, 0744)
			if err != nil {
				return err
			}
		}

		log.Println("[ ] Copy " + srcPath + " into " + destPath)

		input, err := ioutil.ReadFile(srcPath)
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(destPath, input, 0744)
		if err != nil {
			return err
		}

		log.Println("[✔] Copy " + srcPath + " into " + destPath)
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

func replacePlaceholders(template, repoName, appName, arch, tag, version string) (str string) {
	str = strings.Replace(template, placeholderForRepoName, repoName, -1)
	str = strings.Replace(str, placeholderForAppName, appName, -1)
	str = strings.Replace(str, placeholderForArch, arch, -1)
	str = strings.Replace(str, placeholderForTag, tag, -1)
	str = strings.Replace(str, placeholderForVersion, version, -1)

	return
}

func replaceSrcDestTemplates(srcFileTemplate, destPathTemplate, repoName, appName, arch, tag, version string) (srcFile string, destPath string) {
	srcFile = replacePlaceholders(srcFileTemplate, repoName, appName, arch, tag, version)
	destPath = replacePlaceholders(destPathTemplate, repoName, appName, arch, tag, version)
	destPath = strings.Replace(destPath, placeholderForSrc, srcFile, -1)

	return
}

func generateDownloadUrl(template, repoName, appName, arch, tag, version, srcFile string) (url string) {
	url = replacePlaceholders(template, repoName, appName, arch, tag, version)
	url = strings.Replace(url, placeholderForSrc, srcFile, -1)

	return
}

func downloadFile(url, destPath string) error {

	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return fmt.Errorf("error on download %s with status code %v", url, response.StatusCode)
	}

	file, err := os.Create(destPath)
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
