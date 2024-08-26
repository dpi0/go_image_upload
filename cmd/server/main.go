package main

import (
	"github.com/dpi0/go_image_upload/internal/config"
	"github.com/dpi0/go_image_upload/internal/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

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
