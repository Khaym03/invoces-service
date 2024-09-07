package main

import (
	"os"
	"path/filepath"

	"github.com/Khaym03/invoces-service/internal/common"
	"github.com/Khaym03/invoces-service/internal/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	// Read in .env on import
	_ "github.com/joho/godotenv/autoload"
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

	port := ":" + common.Getenv("PORT", "2003")
	app.Listen(port)
}

func faviconMiddleware(c *fiber.Ctx) error {
	if c.Path() == "/favicon.ico" {
		return c.Status(fiber.StatusNotFound).SendString("No favicon available")
	}
	return c.Next()
}
