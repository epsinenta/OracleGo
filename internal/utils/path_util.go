package utils

import (
	"os"
	"strings"
)

func GetPath(relativePath string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	dirs := strings.Split(wd, "/")
	var path string
	for i := len(dirs) - 1; i >= 0; i-- {
		if dirs[i] == "OracleGo" {
			path = strings.Join(dirs[:i+1], "/")
			break
		}
	}

	return path + "/" + relativePath, nil
}
