package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Khaym03/invoces-service/internal/core/service/pdfinvoice"
	"github.com/Khaym03/invoces-service/internal/models"
)

func GenerateInvoiceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Solo se permiten solicitudes POST", http.StatusBadRequest)
		return
	}

	ig := pdfinvoice.Service()

	// Crear un nombre de archivo único para el HTML
	htmlFilename := fmt.Sprintf("invoice-%d.html", time.Now().UnixNano())
	htmlFilePath := filepath.Join("generated", htmlFilename)

	invoice, err := decodeInvoiceRequest(r)
	if err != nil {
		http.Error(w, "Error al decodificar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	// Guardar el archivo HTML
	err = os.MkdirAll(filepath.Dir(htmlFilePath), os.ModePerm)
	if err != nil {
		http.Error(w, "Error al crear el directorio", http.StatusInternalServerError)
		return
	}

	htmlContent := ig.BuildHTML(*invoice)

	err = os.WriteFile(htmlFilePath, htmlContent.Bytes(), 0644)
	if err != nil {
		http.Error(w, "Error al escribir el archivo", http.StatusInternalServerError)
		return
	}

	go ig.GeneratePDF(htmlFilePath)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Archivo HTML generado /%s\n", htmlFilename)
}

func decodeInvoiceRequest(r *http.Request) (*models.InvoiceInput, error) {
	var invoice models.InvoiceInput

	err := json.NewDecoder(r.Body).Decode(&invoice)
	if err != nil {
		return nil, err
	}

	return &invoice, nil
}
