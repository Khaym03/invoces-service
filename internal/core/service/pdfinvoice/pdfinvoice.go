package pdfinvoice

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/Khaym03/invoces-service/internal/common"
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

type InvoiceGenerator struct {
	ctx    context.Context
	cancel context.CancelFunc
	queue  chan string
	wg     sync.WaitGroup
}

func (ig *InvoiceGenerator) InitContext() {
	if ig.ctx == nil {
		opts := append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("headless", true),
			chromedp.Flag("disable-gpu", true),
			chromedp.Flag("no-sandbox", true),
		)

		execAllocator, _ := chromedp.NewExecAllocator(context.Background(), opts...)
		ig.ctx, ig.cancel = chromedp.NewContext(execAllocator)
		ig.queue = make(chan string, 100) // Tamaño de la cola
		ig.startWorker()
	}
}

func (ig *InvoiceGenerator) CloseContext() {
	if ig.cancel != nil {
		ig.cancel()
		close(ig.queue)
		ig.wg.Wait()
	}
}

func (ig *InvoiceGenerator) startWorker() {
	ig.wg.Add(1)
	go func() {
		defer ig.wg.Done()
		for filename := range ig.queue {
			ig.GeneratePDF(filename, time.Now())
		}
	}()
}

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

func (ig *InvoiceGenerator) GeneratePDF(filename string, slapse time.Time) {
	urlFilename := filepath.Base(filename)
	port := common.Getenv("PORT", "2003")
	url := fmt.Sprintf("http://localhost:%s/html-templates/%s", port, urlFilename)
	fmt.Println(url)

	var buf []byte

	err := chromedp.Run(ig.ctx, ig.GeneratePDFFromURL(url, &buf))
	if err != nil {
		log.Printf("Error generating PDF: %v", err)
		return
	}

	err = ig.SavePDFToFile(filename, buf)
	if err != nil {
		log.Printf("Error saving PDF: %v", err)
		return
	}

	fmt.Println(time.Since(slapse))
}

// Guardar el PDF en el directorio predeterminado
func (ig *InvoiceGenerator) SavePDFToFile(filename string, pdfData []byte) error {
	pdfName := fmt.Sprintf("invoice-%d.pdf", time.Now().UnixNano())
	pdfFilename := filepath.Join(
		"invoices", pdfName,
	)

	fmt.Println("PDF guardardo ", pdfName)

	return os.WriteFile(pdfFilename, pdfData, 0644)
}

func (ig *InvoiceGenerator) EnqueuePDFGeneration(filename string) {
	ig.queue <- filename
}
