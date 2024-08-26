package storage

import (
	"io"
	"os"
	"path/filepath"
)

const UploadDir = "./uploads"

// SaveFile saves the uploaded file to the specified path
func SaveFile(src io.Reader, dstPath string) error {
	// Create the destination file
	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy the content from the uploaded file to the destination file
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return nil
}

// GetFilePath returns the file path based on the UUID and file name
func GetFilePath(id, name string) (string, error) {
	filePath := filepath.Join(UploadDir, id+"_"+name)

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", err
	}

	return filePath, nil
}

// DeleteFile deletes the file at the specified path
func DeleteFile(filePath string) error {
	if err := os.Remove(filePath); err != nil {
		return err
	}

	return nil
}
