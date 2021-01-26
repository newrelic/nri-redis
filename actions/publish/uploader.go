package main


func main() {
	// get version from env
	version := "0.0.0"
	// take schema
	fileContent, err := readFileContent("filepath")

	if err != nil{
		// do something
	}

	//handle err
	// parse schema
	config := parseConfig(fileContent)
	// iterate over config
	uploadArtifacts(config, version)
	// --> replace variables
	// --> copy (move?) file
}

type uploadArtifactConfig struct {
	src string `json:"src"`
	dest string `json:"dest"`
	arch []string `json:"arch"`
}

type uploadConfig []uploadArtifactConfig

func readFileContent(filePath string) (string, error){
	return "", nil
}

func parseConfig(fileContent string) uploadConfig {
	return uploadConfig{}
}

func uploadArtifact(config uploadArtifactConfig){}

func uploadArtifacts(config uploadConfig, version string){}

func replacePlaceholders(srcTemplate, destTemplate, arch, binaryName, version string) (string, string){

	srcFileName := ""
	destPath := ""
	return srcFileName, destPath
}