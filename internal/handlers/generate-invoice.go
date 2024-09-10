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

type handler struct {
	InvoiceService *pdfinvoice.InvoiceGenerator
}

func Handler(pdfSrv *pdfinvoice.InvoiceGenerator) *handler {
	return &handler{
		InvoiceService: pdfSrv,
	}
}

func (h *handler) GenerateInvoiceHandler(c *fiber.Ctx) error {
	// Crear un nombre de archivo único para el HTML
	fileName := fmt.Sprintf("invoice-%d", time.Now().UnixNano())
	htmlFilename := fileName + ".html"
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

	htmlContent := h.InvoiceService.BuildHTML(*invoice)

	err = os.WriteFile(htmlFilePath, htmlContent.Bytes(), 0644)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error al escribir el archivo")
	}

	filename := h.InvoiceService.WaitForPDFGeneration(fileName)
	fmt.Printf("PDF generado: %s\n", filename)

	return c.Status(fiber.StatusOK).SendString(filename)
}

func decodeInvoiceRequest(c *fiber.Ctx) (*models.InvoiceInput, error) {
	var invoice models.InvoiceInput
	err := c.BodyParser(&invoice)
	return &invoice, err
}
