package pdfinvoice

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	adapters "github.com/Khaym03/invoces-service/internal/adapters/pdfstorage"
	"github.com/Khaym03/invoces-service/internal/common"
	"github.com/Khaym03/invoces-service/internal/components"
	"github.com/Khaym03/invoces-service/internal/core/ports"
	"github.com/Khaym03/invoces-service/internal/core/service/emailsender"
	"github.com/Khaym03/invoces-service/internal/models"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

var intance *InvoiceGenerator
var once sync.Once

func Service() *InvoiceGenerator {
	once.Do(func() {
		intance = &InvoiceGenerator{
			Email:      emailsender.Service(),
			PDFStorage: adapters.NewStoreLocally(),
		}
	})
	return intance
}

type InvoiceGenerator struct {
	ctx        context.Context
	cancel     context.CancelFunc
	queue      chan string
	done       chan string
	wg         sync.WaitGroup
	Email      ports.EmailSender
	PDFStorage ports.PDFStorage
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
		ig.done = make(chan string, 100)
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
			pdfUrl := ig.GeneratePDF(filename, time.Now())
			ig.done <- pdfUrl
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

func (ig *InvoiceGenerator) GeneratePDF(fileName string, slapse time.Time) string {
	htmlFile := fileName + ".html"
	port := common.Getenv("PORT", "2003")
	url := fmt.Sprintf("http://localhost:%s/html-templates/%s", port, htmlFile)
	pdfName := fileName + ".pdf"
	fmt.Println(url)

	var buf []byte

	err := chromedp.Run(ig.ctx, ig.GeneratePDFFromURL(url, &buf))
	if err != nil {
		log.Printf("Error generating PDF: %v", err)
		return ""
	}

	pdfUrl, err := ig.PDFStorage.Save(pdfName, buf)
	if err != nil {
		log.Printf("Error saving PDF: %v", err)
		return ""
	}

	fmt.Println(time.Since(slapse))

	return pdfUrl
}

func (ig *InvoiceGenerator) EnqueuePDFGeneration(filename string) {
	ig.queue <- filename
}

func (ig *InvoiceGenerator) WaitForPDFGeneration(filename string) string {
	ig.EnqueuePDFGeneration(filename)
	return <-ig.done
}
