package pdfinvoice

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/Khaym03/invoces-service/internal/components"
	"github.com/Khaym03/invoces-service/internal/models"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

var intance *InvoiceGenerator
var once sync.Once

func Service() *InvoiceGenerator {
	once.Do(func() {
		intance = &InvoiceGenerator{}
	})
	return intance
}

type InvoiceGenerator struct{}

func (ig *InvoiceGenerator) GeneratePDFFromURL(urlstr string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.ActionFunc(func(ctx context.Context) error {
			params := page.PrintToPDF().
				WithMarginTop(0).
				WithMarginBottom(0).
				WithMarginLeft(0).
				WithMarginRight(0).
				WithPrintBackground(true)

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

func (ig *InvoiceGenerator) BuildHTML(i models.InvoiceInput) bytes.Buffer {
	var buf bytes.Buffer
	root := components.Root(i)

	root.Render(context.Background(), &buf)

	return buf
}

func (ig *InvoiceGenerator) GeneratePDF(filename string) {
	// Eliminar el directorio 'generated/' del nombre del archivo
	urlFilename := filepath.Base(filename)
	url := fmt.Sprintf("http://localhost:3000/html-templates/%s", urlFilename)
	fmt.Println(url)

	var buf []byte
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	err := chromedp.Run(ctx, ig.GeneratePDFFromURL(url, &buf))
	if err != nil {
		fmt.Println("Error al generar el PDF:", err)
		return
	}

	err = ig.SavePDFToFile(filename, buf)
	if err != nil {
		fmt.Printf("Error al escribir el PDF en 'invoices': %v", err)
		return
	}

	fmt.Printf("PDF generado con éxito\n")
}

// Guardar el PDF en el directorio predeterminado
func (ig *InvoiceGenerator) SavePDFToFile(filename string, pdfData []byte) error {
	pdfFilename := filepath.Join(
		"invoices", fmt.Sprintf("invoice-%d.pdf", time.Now().UnixNano()),
	)

	return os.WriteFile(pdfFilename, pdfData, 0644)
}
