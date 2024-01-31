package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/fingerprintjs/fingerprint-pro-server-api-go-sdk/v4/config"
)

var files = []string{"README.md", "docs", ".swagger-codegen"}
var filesToKeep = []string{"docs/DecryptionKey.md", "docs/SealedResults.md"}
var pathPrefix = "sdk"

func main() {
	moveFilesToKeepToTmpDir()
	bumpConfigVersion()
	generateSwagger()
	moveFiles()
	getExamples()
	fixFingerPrintApiMdFile()
	moveFilesToKeepFromTmpDir()
	formatCode()
}

func ensureTmpDir(paths ...string) {
	fullPath := fmt.Sprintf("tmp/%s", strings.Join(paths, "/"))

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		err := os.Mkdir(fullPath, 0755)

		if err != nil {
			log.Fatal(err)
		}
	}
}

func moveFilesToKeepToTmpDir() {
	ensureTmpDir()
	ensureTmpDir("docs")

	for _, file := range filesToKeep {
		filePath := fmt.Sprintf("%s", file)
		newFilePath := fmt.Sprintf("tmp/%s", file)

		err := os.Rename(filePath, newFilePath)

		if err != nil {
			log.Fatal(err)
		}
	}
}

func moveFilesToKeepFromTmpDir() {
	for _, file := range filesToKeep {
		filePath := fmt.Sprintf("./tmp/%s", file)
		newFilePath := fmt.Sprintf("%s", file)

		err := os.Rename(filePath, newFilePath)

		if err != nil {
			log.Fatal(err)
		}
	}
}

func getVersion() string {
	var version string

	envVersion := os.Getenv("VERSION")

	if envVersion != "" {
		version = envVersion
	} else {
		configFile := config.ReadConfig("./config.json")
		version = configFile.PackageVersion
	}

	return version
}

func bumpConfigVersion() {
	version := getVersion()

	configFile := config.ReadConfig("./config.json")

	if configFile.PackageVersion == version {
		log.Println("Version is up to date")
		return
	}

	configFile.PackageVersion = version

	configContents, err := json.MarshalIndent(configFile, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	if err = os.WriteFile("./config.json", configContents, 0644); err != nil {
		log.Fatal(err)
	}
}

func removeFileOrDirIfExists(path string) {
	if stat, err := os.Stat(path); err == nil {
		var err error

		if stat.IsDir() {
			err = os.RemoveAll(path)
		} else {
			err = os.Remove(path)
		}

		if err != nil {
			log.Fatal(err)
		}
	}
}

func cleanupOldFiles() {
	for _, filePath := range files {
		removeFileOrDirIfExists(filePath)
	}
}

func moveFiles() {
	cleanupOldFiles()

	for _, file := range files {

		filePath := fmt.Sprintf("%s/%s", pathPrefix, file)
		newFilePath := fmt.Sprintf("./%s", file)

		err := os.Rename(filePath, newFilePath)

		if err != nil {
			log.Fatal(err)
		}
	}
}

func generateSwagger() {
	cmd := exec.Command(
		"java",
		"-jar",
		"./bin/swagger-codegen-cli.jar",
		"generate",
		"-t",
		"./template",
		"-l",
		"go",
		"-i",
		"res/fingerprint-server-api.yaml",
		"-o",
		"./sdk",
		"-c",
		"config.json")

	out, cmdErr := cmd.Output()

	if cmdErr != nil {
		log.Fatal(cmdErr)
	}

	fmt.Println(string(out))
}

func getExamples() {
	list := []string{
		"visits_limit_1.json",
		"visits_limit_500.json",
		"webhook.json",
		"get_event.json",
		"get_event_all_errors.json",
		"get_event_extra_fields.json",
		"get_event_403_error.json",
		"get_event_404_error.json",
		"get_event_botd_failed_error.json",
		"get_event_botd_too_many_requests_error.json",
		"get_event_identification_failed_error.json",
		"get_event_identification_too_many_requests_error.json",
		"visits_too_many_requests_error.json",
	}

	for _, file := range list {
		cmd := exec.Command("curl", "-o", fmt.Sprintf("./test/mocks/%s", file), fmt.Sprintf("https://fingerprintjs.github.io/fingerprint-pro-server-api-openapi/examples/%s", file))
		_, err := cmd.Output()

		if err != nil {
			log.Fatal(err)
		}
	}

}

/**
 * Fixes a bug with generated file in "docs/FingerprintApi.md" which contains invalid title generated by swagger
 */
func fixFingerPrintApiMdFile() {
	token := "{{classname}}"
	target := "FingerprintApi"
	targetsToRemove := []string{"**optional** | ***FingerprintApiGetVisitsOpts** | optional parameters | nil if no parameters"}
	filePath := "docs/FingerprintApi.md"

	fileContents, err := os.ReadFile(filePath)

	if err != nil {
		log.Fatal(err)
	}

	fileContents = []byte(strings.Replace(string(fileContents), token, target, -1))
	fileContentsArray := strings.Split(string(fileContents), "\n")
	var fileContentsArrayResult []string

	// Fixes markdown table for optional parameters, by default swagger-codegen applies new line there which breaks the table.
	for _, line := range fileContentsArray {
		for _, targetToRemove := range targetsToRemove {
			if line != targetToRemove {
				fileContentsArrayResult = append(fileContentsArrayResult, strings.Replace(line, targetToRemove, "", -1))
			}
		}
	}

	err = os.WriteFile(filePath, []byte(strings.Join(fileContentsArrayResult, "\n")), 0644)

	if err != nil {
		log.Fatal(err)
	}
}

func formatCode() {
	cmd := exec.Command("go", "fmt")
	cmd.Dir = "./sdk"

	_, err := cmd.Output()

	if err != nil {
		log.Fatal(err)
	}
}
