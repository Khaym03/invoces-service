package main

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	"log"
	// "os"

	"github.com/Khaym03/invoces-service/internal/components"
	"github.com/Khaym03/invoces-service/internal/models"

	// "github.com/chromedp/cdproto/page"
	// "github.com/chromedp/chromedp"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

// var pdf = template.Must(template.ParseFiles("templates/test.html"))

type Invoice struct {
	Name string
}

func PdfTest() string {
	var buf bytes.Buffer

	description := models.InvoiceDescription{
		Id:       123123,
		Date:     time.Now(),
		DueDate:  time.Now(),
		TotalDue: 50.00,
	}

	// des := components.InvoiceDescription(description)

	root := components.Root(description)

	root.Render(context.Background(), &buf)

	return buf.String()
}

func main() {
	ExampleNewPDFGenerator()
}

func ExampleNewPDFGenerator() {

	// Create new PDF generator
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		log.Fatal(err)
	}

	// Set global options
	pdfg.Dpi.Set(300)
	pdfg.Orientation.Set(wkhtmltopdf.OrientationPortrait)
	pdfg.Grayscale.Set(false)

	htmlPdf := PdfTest()

	fmt.Println(htmlPdf)

	pdfg.AddPage(wkhtmltopdf.NewPageReader(strings.NewReader(htmlPdf)))

	pdfg.TOC.FooterRight.Set("[page]")
	pdfg.TOC.FooterFontSize.Set(10)
	pdfg.TOC.Zoom.Set(0.95)

	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)

	const margin = 16
	pdfg.MarginTop.Set(margin)
	pdfg.MarginLeft.Set(margin)
	pdfg.MarginBottom.Set(margin)
	pdfg.MarginLeft.Set(margin)

	// Create PDF document in internal buffer
	err = pdfg.Create()
	if err != nil {
		log.Fatal(err)
	}

	// Write buffer contents to file on disk
	err = pdfg.WriteFile("./simplesample.pdf")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Done")
	// Output: Done
}
