package file

import (
	"fmt"
	"io/ioutil"
	"os"
)

func ReadLocalFile() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %v", err)
	}
	fileContent, err := ioutil.ReadFile(fmt.Sprintf("%v/%v", homeDir, ".kube/config"))
	if err != nil {
		return "", err
	}
	return string(fileContent), nil
}
