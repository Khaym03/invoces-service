package handlers

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Khaym03/invoces-service/internal/core/service/pdfinvoice"
	"github.com/Khaym03/invoces-service/internal/models"
	"github.com/gofiber/fiber/v2"
)

func GenerateInvoiceHandler(c *fiber.Ctx) error {
	ig := pdfinvoice.Service()

	// Crear un nombre de archivo único para el HTML
	htmlFilename := fmt.Sprintf("invoice-%d.html", time.Now().UnixNano())
	htmlFilePath := filepath.Join("html-templates", htmlFilename)

	invoice, err := decodeInvoiceRequest(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Error al decodificar el cuerpo de la solicitud")
	}

	// Guardar el archivo HTML
	err = os.MkdirAll(filepath.Dir(htmlFilePath), os.ModePerm)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error al crear el directorio")
	}

	htmlContent := ig.BuildHTML(*invoice)

	err = os.WriteFile(htmlFilePath, htmlContent.Bytes(), 0644)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error al escribir el archivo")
	}

	go ig.GeneratePDF(htmlFilePath)

	fmt.Printf("Archivo HTML generado /%s\n", htmlFilename)

	return c.Status(fiber.StatusOK).SendString(fmt.Sprintf("Archivo HTML generado /%s\n", htmlFilename))
}

func decodeInvoiceRequest(c *fiber.Ctx) (*models.InvoiceInput, error) {
	var invoice models.InvoiceInput
	err := c.BodyParser(&invoice)
	return &invoice, err
}
