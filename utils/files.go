package utils

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

type FileInfo struct {
	Location  string
	Extension string
}

func SaveMultipartFile(file multipart.File, fileInfo *multipart.FileHeader, dest string, allowed []string) (FileInfo, int, error) {
	var fileExtension, fileLocation string
	var err error
	var result FileInfo
	var notAllowed = true

	// start check extension
	splitContenttype := strings.Split(fileInfo.Filename, ".")
	fileExtension = splitContenttype[len(splitContenttype)-1]
	for _, s := range allowed {
		if strings.ToLower(s) == fileExtension {
			notAllowed = false
		}
	}
	if notAllowed {
		return result, http.StatusBadRequest, errors.New("file is not support")
	}
	fileLocation = dest + "." + fileExtension
	// end check extension

	// start upload
	fileTarget, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		err = errors.New("something error on server")
		return FileInfo{}, http.StatusInternalServerError, err
	}
	defer func(fileTarget *os.File) {
		err := fileTarget.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(fileTarget)
	if _, err := io.Copy(fileTarget, file); err != nil {
		err = errors.New("something error on server")
		return FileInfo{}, http.StatusInternalServerError, err
	}
	// end upload

	// set file information
	result.Extension = fileExtension
	result.Location = fileLocation

	return result, http.StatusCreated, nil
}
