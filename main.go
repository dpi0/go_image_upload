package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const uploadDir = "./uploads"

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Create upload directory if it doesn't exist
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.Mkdir(uploadDir, 0755)
	}

	// Routes
	e.POST("/upload", uploadFile)
	e.GET("/file/:id/:name", downloadFile)
	e.DELETE("/file/:id/:name", deleteFile)
	e.GET("/files", listFiles)

	// Start server
	e.Logger.Fatal(e.Start(":8081"))
}

// uploadFile handles file uploads
func uploadFile(c echo.Context) error {
	// Get the uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		return c.String(http.StatusBadRequest, "Failed to get file")
	}

	// Generate a short UUID (first 5 characters)
	shortUUID := uuid.New().String()[:5]

	// Save the file with the short UUID as part of the name
	dst, err := os.Create(filepath.Join(uploadDir, shortUUID+"_"+file.Filename))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to save file")
	}
	defer dst.Close()

	src, err := file.Open()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to open file")
	}
	defer src.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return c.String(http.StatusInternalServerError, "Failed to copy file")
	}

	// URL encode the file name to handle spaces and special characters
	encodedFileName := url.PathEscape(file.Filename)

	// Return the download URL
	url := fmt.Sprintf("http://%s/file/%s/%s", c.Request().Host, shortUUID, encodedFileName)
	return c.JSON(http.StatusOK, map[string]string{
		"url":  url,
		"name": file.Filename,
	})
}

// downloadFile handles file download requests
func downloadFile(c echo.Context) error {
	id := c.Param("id")
	name := c.Param("name")

	// URL decode the file name to get the original file name
	decodedFileName, err := url.PathUnescape(name)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid file name")
	}

	// Search for the file with the corresponding UUID and filename
	filePath := filepath.Join(uploadDir, id+"_"+decodedFileName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return c.String(http.StatusNotFound, "File not found")
	}

	// Serve the file
	return c.File(filePath)
}

// deleteFile handles file deletion requests
func deleteFile(c echo.Context) error {
	id := c.Param("id")
	name := c.Param("name")

	// URL decode the file name to get the original file name
	decodedFileName, err := url.PathUnescape(name)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid file name")
	}

	// Search for the file with the corresponding UUID and filename
	filePath := filepath.Join(uploadDir, id+"_"+decodedFileName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return c.String(http.StatusNotFound, "File not found")
	}

	// Delete the file
	if err := os.Remove(filePath); err != nil {
		return c.String(http.StatusInternalServerError, "Failed to delete file")
	}

	return c.String(http.StatusOK, "File deleted successfully")
}

// listFiles returns a list of all uploaded files in JSON format
func listFiles(c echo.Context) error {
	var files []map[string]string

	err := filepath.Walk(uploadDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
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

				fileURL := fmt.Sprintf("http://%s/file/%s/%s", c.Request().Host, shortUUID, encodedFileName)
				files = append(files, map[string]string{
					"url":  fileURL,
					"name": originalFileName,
				})
			}
		}
		return nil
	})

	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to list files")
	}

	return c.JSON(http.StatusOK, files)
}
