package utils

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

func SavePicture(file multipart.File, header *multipart.FileHeader) (string, error) {
	os.MkdirAll("uploads", os.ModePerm)
	filename := filepath.Join("uploads", header.Filename)
	out, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	return filename, err
}
