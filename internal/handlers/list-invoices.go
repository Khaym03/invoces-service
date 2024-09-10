package handlers

import (
	"os"

	"github.com/gofiber/fiber/v2"
)

func (h *handler) ListInvoices(c *fiber.Ctx) error {
	entries, err := os.ReadDir("invoices")
	if err != nil {
		c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	var fileNames []string
	for _, entry := range entries {
		fn := entry.Name()
		if fn != ".gitkeep" {
			fileNames = append(fileNames, fn)
		}
	}

	c.JSON(fileNames)

	return nil
}
