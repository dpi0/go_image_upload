package main

import (
	"os"

	"github.com/dpi0/go_image_upload/internal/config"
	"github.com/dpi0/go_image_upload/internal/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

}

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Initialize configuration
	config.InitConfig()

	// Register routes
	handlers.RegisterRoutes(e)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
