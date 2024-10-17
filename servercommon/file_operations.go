package servercommon

import (
	"errors"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
)

const (
	FILE_NAME_LEN int = 10
)

func generateRandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func IsFileExists(fileName string) (bool, error) {
	_, err := os.Stat(fileName)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return true, err
}

func GetUniqFileName(folder string, extension string) (string, error) {
	path := ""
	for {
		path = filepath.Join(folder, generateRandomString(FILE_NAME_LEN))
		if extension != "" {
			path = path + extension
		}
		isExist, err := IsFileExists(path)
		if err != nil {
			return "", err
		}
		if !isExist {
			break
		}
	}
	return path, nil
}

func GetBaseFolder() string {
	basePath := os.Getenv(BASE_FOLDER_ENV)
	if basePath != "" {
		return basePath
	}
	return BASE_FOLDER
}

func CreateFolder(path string) error {
	err := os.MkdirAll(path, os.ModeDir)
	if err == nil || os.IsExist(err) {
		return nil
	}
	return err
}

func GetProjectFolder(projectId uint) string {
	baseFolder := GetBaseFolder()
	baseFolder = filepath.Join(baseFolder, PROJECTS_FOLDER, strconv.FormatUint(uint64(projectId), 10))
	return baseFolder
}
