package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Khaym03/invoces-service/internal/components"
	"github.com/Khaym03/invoces-service/internal/models"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func main() {
	http.HandleFunc("/assets/css/output.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		http.ServeFile(w, r, "./assets/css/output.css")

	})

	dir, _ := os.Getwd()

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(filepath.Join(dir, "assets")))))

	http.HandleFunc("/generate", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Solo se permiten solicitudes POST", http.StatusBadRequest)
			return
		}

		var invoice models.InvoiceDescription
		err := json.NewDecoder(r.Body).Decode(&invoice)
		if err != nil {
			http.Error(w, "Error al decodificar el cuerpo de la solicitud", http.StatusBadRequest)
			return
		}

		// Crear un nombre de archivo único para el HTML
		filename := fmt.Sprintf("invoice-%d.html", time.Now().UnixNano())
		folder := filepath.Join("generated", filename)

		// Guardar el archivo HTML
		err = os.MkdirAll(filepath.Dir(folder), os.ModePerm)
		if err != nil {
			http.Error(w, "Error al crear el directorio", http.StatusInternalServerError)
			return
		}

		_, buf := PdfTest()

		err = os.WriteFile(folder, buf.Bytes(), 0644)
		if err != nil {
			http.Error(w, "Error al escribir el archivo", http.StatusInternalServerError)
			return
		}

		// Servir el archivo HTML generado
		http.HandleFunc("/"+filename, func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, folder)
		})

		go GeneratePDF(folder)

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Archivo HTML generado y servido en /%s\n", filename)
	})

	fmt.Println("Listening on :3000")
	http.ListenAndServe(":3000", nil)

}

// Función para generar el PDF
func PrintToPDF(urlstr string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.ActionFunc(func(ctx context.Context) error {
			params := page.PrintToPDF().
				WithMarginTop(0).
				WithMarginBottom(0).
				WithMarginLeft(0).
				WithMarginRight(0).
				WithPrintBackground(true)

			// Generar el PDF con fondo transparente (opcional)
			buf, _, err := params.Do(ctx)
			if err != nil {
				return err
			}

			// Asignar el contenido del PDF al buffer
			*res = buf
			return nil
		}),
	}
}

// Función para generar el PDF a partir de contenido HTML dinámico
func PrintToPDFFromHTML(htmlContent string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate("data:text/html," + htmlContent),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// Generar el PDF con fondo transparente (opcional)
			buf, _, err := page.PrintToPDF().WithPrintBackground(true).Do(ctx)
			if err != nil {
				return err
			}

			// Asignar el contenido del PDF al buffer
			*res = buf
			return nil
		}),
	}
}

func PdfTest() (string, bytes.Buffer) {
	var buf bytes.Buffer

	description := models.InvoiceDescription{
		Id:       123123,
		Date:     time.Now(),
		DueDate:  time.Now(),
		TotalDue: 50.00,
	}

	root := components.Root(description)

	root.Render(context.Background(), &buf)

	return buf.String(), buf
}

func PDFFromBrowser() {
	// Crear un contexto
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// URL de la página a convertir a PDF
	url := "http://localhost:3000/x"

	// Buffer para almacenar el contenido del PDF
	var buf []byte

	// Ejecutar la tarea de generar el PDF
	if err := chromedp.Run(ctx, PrintToPDF(url, &buf)); err != nil {
		log.Fatal(err)
	}

	// Escribir el contenido del PDF a un archivo
	// if err := os.WriteFile("sample.pdf", buf, 0644); err != nil {
	// 	log.Fatal(err)
	// }

	// Ejemplo de uso
	err := writeToFile("browser.pdf", buf, 0644)
	if err != nil {
		log.Fatalf("Error al escribir archivo: %v", err)
	}

	fmt.Println("Se ha escrito sample.pdf")
}

func writeToFile(filePath string, content []byte, perm os.FileMode) error {
	// Obtener el directorio actual
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error obteniendo ruta absoluta: %v", err)
	}

	// Construir la ruta completa del archivo
	fullPath := filepath.Join(currentDir, "generated-pdf", filepath.Base(filePath))

	// Escribir el contenido del PDF al archivo
	err = os.WriteFile(fullPath, content, perm)
	if err != nil {
		return fmt.Errorf("error escribiendo archivo %s: %v", fullPath, err)
	}

	return nil
}

func GeneratePDF(filename string) {
	// Eliminar el directorio 'generated/' del nombre del archivo
	urlFilename := filepath.Base(filename)
	url := fmt.Sprintf("http://localhost:3000/%s", urlFilename)
	fmt.Println(url)

	var buf []byte
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	err := chromedp.Run(ctx, PrintToPDF(url, &buf))
	if err != nil {
		log.Println("Error al generar el PDF:", err)
		return
	}

	// Guardar el PDF en el directorio predeterminado
	pdfFilename := filepath.Join("generated-pdf", fmt.Sprintf("invoice-%d.pdf", time.Now().UnixNano()))
	err = os.WriteFile(pdfFilename, buf, 0644)
	if err != nil {
		log.Printf("Error al escribir el PDF en %s: %v", "generated-pdf", err)
		return
	}

	fmt.Printf("PDF generado con éxito: %s\n", pdfFilename)
}
