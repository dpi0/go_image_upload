package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

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

	// Return the download URL
	url := fmt.Sprintf("http://%s/file/%s/%s", c.Request().Host, shortUUID, file.Filename)
	return c.JSON(http.StatusOK, map[string]string{
		"url":  url,
		"name": file.Filename,
	})
}

// downloadFile handles file download requests
func downloadFile(c echo.Context) error {
	id := c.Param("id")
	name := c.Param("name")

	// Search for the file with the corresponding UUID and filename
	filePath := filepath.Join(uploadDir, id+"_"+name)
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

	// Search for the file with the corresponding UUID and filename
	filePath := filepath.Join(uploadDir, id+"_"+name)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return c.String(http.StatusNotFound, "File not found")
	}

	// Delete the file
	if err := os.Remove(filePath); err != nil {
		return c.String(http.StatusInternalServerError, "Failed to delete file")
	}

	return c.String(http.StatusOK, "File deleted successfully")
}
