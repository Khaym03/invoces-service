package ports

import "github.com/Khaym03/invoces-service/internal/models"

type EmailSender interface {
	Send(to models.CustomerDetails, pdfURL string) error
}

type PDFStorage interface {
	Save(pdfName string, pdfData []byte) (string, error)
}
