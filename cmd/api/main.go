package main

import (
	"os"
	"path/filepath"

	"github.com/Khaym03/invoces-service/internal/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New())
	app.Use(faviconMiddleware)

	rootDir, _ := os.Getwd()

	app.Static("/assets", filepath.Join(rootDir, "assets"))
	app.Static("/invoices", filepath.Join(rootDir, "invoices"))
	app.Static("/html-templates", filepath.Join(rootDir, "html-templates"))

	app.Post("/generate-invoice", handlers.GenerateInvoiceHandler)

	app.Listen(":3000")
}

func faviconMiddleware(c *fiber.Ctx) error {
	if c.Path() == "/favicon.ico" {
		return c.Status(fiber.StatusNotFound).SendString("No favicon available")
	}
	return c.Next()
}
