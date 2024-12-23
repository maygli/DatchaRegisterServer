package servercommon

import (
	"errors"
	"io"
	"log/slog"
	"math/rand"
	"net/http"
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

func GetProfileFolder(profileId uint) string {
	baseFolder := GetBaseFolder()
	baseFolder = filepath.Join(baseFolder, PROFILE_FOLDER, strconv.FormatUint(uint64(profileId), 10))
	return baseFolder
}

// Remove previous file and save file from multipart/form data.
//
//			r - received request
//			key - file key in the form data
//			folder - base folder
//			fileName - base file name without extension. Extension will get from form form file name.
//		            random path will be added to file name to avoid cache by browser or frontend
//		            reactivity system.
//	   prevFileName - previous file name which will be removed
//
// This function returns file path in case of success or err in case of error
func ReplaceFormFile(r *http.Request, key string, folder string, fileName string, prevFilePath string) (string, error) {
	if prevFilePath != "" {
		os.Remove(prevFilePath)
	}
	return SaveFormFile(r, key, folder, fileName)
}

// Save form file to the disk.
//
//		r - received request
//		key - file key in the form data
//		folder - base folder
//		fileName - base file name without extension. Extension will get from form form file name.
//	            random path will be added to file name to avoid cache by browser or frontend
//	            reactivity system.
//
// This function returns file path in case of success or err in case of error
func SaveFormFile(r *http.Request, key string, folder string, fileName string) (string, error) {
	file, header, err := r.FormFile(key)
	if err == nil {
		defer file.Close()
		err := CreateFolder(folder)
		if err != nil {
			// Do nothing in case it is impossible to create project
			return "", err
		}
		rand := generateRandomString(8)
		imageFileName := filepath.Join(folder, fileName+"-"+rand+filepath.Ext(header.Filename))
		// Copy the file data to buffer
		f, err := os.OpenFile(imageFileName, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			return "", err
		}
		defer f.Close()
		_, err = io.Copy(f, file)
		if err != nil {
			return "", err
		}
		return imageFileName, nil
	}
	// In case of we don't read form file - don't return error. But report error to log
	slog.Error("can't read form file. Error: " + err.Error())
	return "", nil
}
