package a3interface

import (
	"fmt"
	"os"
)

// GetArmaDir returns the Arma 3 executable directory. It will not account for symlinks.
func GetArmaDir() (string, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error getting executable directory: %s", err.Error())
	}
	return workingDir, nil
}
