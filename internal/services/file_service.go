package services

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/dpi0/go_image_upload/internal/storage"
	"github.com/dpi0/go_image_upload/pkg/utils"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// UploadFile handles the business logic for uploading a file
func UploadFile(c echo.Context) (string, error) {
	// Get the uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		log.Error().Err(err).Msg("Failed to get file")
		return "", fmt.Errorf("failed to get file: %w", err)
	}

	// Get the original file name without UUID prefix
	originalFileName := file.Filename

	// Check if a file with the same name already exists
	files, err := os.ReadDir(storage.UploadDir)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read directory")
		return "", fmt.Errorf("failed to read directory: %w", err)
	}

	for _, f := range files {
		// Remove UUID prefix from file name for comparison
		filenameWithoutUUID := strings.Split(f.Name(), "_")[1]
		if filenameWithoutUUID == originalFileName {
			log.Warn().Msgf("File %s already exists. Upload canceled.", originalFileName)
			// Return a JSON error response
			return `{"error": "file already exists"}`, fmt.Errorf("file %s already exists", originalFileName)
		}
	}
	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		log.Error().Err(err).Msg("Failed to open file")
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// Generate a short UUID (first 5 characters)
	shortUUID := utils.GenerateShortUUID()

	// Destination file path
	dstPath := filepath.Join(storage.UploadDir, shortUUID+"_"+file.Filename)

	// Save the file to the storage
	if err := storage.SaveFile(src, dstPath); err != nil {
		log.Error().Err(err).Msg("Failed to save file")
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	// Return the download URL
	url := fmt.Sprintf("http://%s/%s/%s", c.Request().Host, shortUUID, file.Filename)
	return url, nil
}

// DownloadFile handles the business logic for downloading a file
func DownloadFile(c echo.Context) (string, error) {
	id := c.Param("id")
	name := c.Param("name")

	// Get the file path from storage
	filePath, err := storage.GetFilePath(id, name)
	if err != nil {
		if os.IsNotExist(err) {
			log.Warn().Err(err).Msg("File not found")
			return "", fmt.Errorf("file not found: %w", err)
		}
		log.Error().Err(err).Msg("Failed to get file path")
		return "", fmt.Errorf("failed to get file path: %w", err)
	}

	return filePath, nil
}

// DeleteFile handles the business logic for deleting a file
func DeleteFile(c echo.Context) error {
	id := c.Param("id")
	name := c.Param("name")

	// Get the file path from storage
	filePath, err := storage.GetFilePath(id, name)
	if err != nil {
		if os.IsNotExist(err) {
			log.Warn().Err(err).Msg("File not found")
			return fmt.Errorf("file not found: %w", err)
		}
		log.Error().Err(err).Msg("Failed to get file path")
		return fmt.Errorf("failed to get file path: %w", err)
	}

	// Delete the file from storage
	if err := storage.DeleteFile(filePath); err != nil {
		log.Error().Err(err).Msg("Failed to delete file")
		return fmt.Errorf("ailed to delete file: %w", err)
	}

	log.Info().Msg("File deleted successfully")
	return nil
}

// ListFiles handles the business logic for listing files
func ListFiles(c echo.Context) ([]map[string]string, error) {
	var files []map[string]string

	err := filepath.Walk(storage.UploadDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Error().Err(err).Msg("Failed to walk directory")
			return err
		}
		if !info.IsDir() {
			fileName := info.Name()
			parts := strings.SplitN(fileName, "_", 2)
			if len(parts) == 2 {
				shortUUID := parts[0]
				originalFileName := parts[1]

				// URL encode the file name to handle spaces and special characters
				encodedFileName := url.PathEscape(originalFileName)

				fileURL := fmt.Sprintf("http://%s/%s/%s", c.Request().Host, shortUUID, encodedFileName)
				files = append(files, map[string]string{
					"url":  fileURL,
					"name": originalFileName,
				})
			}
		}
		return nil
	})

	if err != nil {
		log.Error().Err(err).Msg("Failed to list files")
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	log.Info().Msg("Files listed successfully")
	return files, nil
}
