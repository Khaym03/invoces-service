package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func main() {
	// Crear un contexto
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// URL de la página a convertir a PDF
	url := "https://www.google.com/"

	// Buffer para almacenar el contenido del PDF
	var buf []byte

	// Ejecutar la tarea de generar el PDF
	if err := chromedp.Run(ctx, printToPDF(url, &buf)); err != nil {
		log.Fatal(err)
	}

	// Escribir el contenido del PDF a un archivo
	if err := os.WriteFile("sample.pdf", buf, 0644); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Se ha escrito sample.pdf")
}

// Función para generar el PDF
func printToPDF(urlstr string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// Generar el PDF con fondo transparente (opcional)
			buf, _, err := page.PrintToPDF().WithPrintBackground(false).Do(ctx)
			if err != nil {
				return err
			}

			// Asignar el contenido del PDF al buffer
			*res = buf
			return nil
		}),
	}
}
