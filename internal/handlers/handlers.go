package handlers

import (
	"github.com/dpi0/go_image_upload/internal/services"

	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func RegisterRoutes(e *echo.Echo) {
	e.POST("/upload", UploadFile)
	e.GET("/file/:id/:name", DownloadFile)
	e.GET("/files", ListFiles)
	e.DELETE("/file/:id/:name", DeleteFile)
}

func UploadFile(c echo.Context) error {

	url, err := services.UploadFile(c)
	if err != nil {
		log.Error().Err(err).Msg("Failed to upload file")
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// Return the download URL in the response
	log.Info().Msg("File uploaded successfully")
	return c.JSON(http.StatusOK, map[string]string{
		"url": url,
	})
}

func DownloadFile(c echo.Context) error {
	filePath, err := services.DownloadFile(c)
	if err != nil {
		if err.Error() == "File not found" {
			log.Warn().Err(err).Msg("File not found")
			return c.String(http.StatusNotFound, err.Error())
		}
		log.Error().Err(err).Msg("Failed to download file")
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// Serve the file for download
	log.Info().Msg("File downloaded successfully")
	return c.File(filePath)
}

func DeleteFile(c echo.Context) error {
	err := services.DeleteFile(c)
	if err != nil {
		if err.Error() == "File not found" {
			log.Warn().Err(err).Msg("File not found")
			return c.String(http.StatusNotFound, err.Error())
		}
		log.Error().Err(err).Msg("Failed to delete file")
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// Return a success message
	log.Info().Msg("File deleted successfully")
	return c.String(http.StatusOK, "File deleted successfully")
}

func ListFiles(c echo.Context) error {
	files, err := services.ListFiles(c)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list files")
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, files)
}
