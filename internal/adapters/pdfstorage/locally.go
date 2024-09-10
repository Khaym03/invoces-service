package adapters

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Khaym03/invoces-service/internal/common"
)

type storeLocally struct{}

func NewStoreLocally() *storeLocally {
	return &storeLocally{}
}

func (sl *storeLocally) Save(pdfName string, pdfData []byte) (string, error) {
	port := common.Getenv("PORT", "2003")
	pdfUrl := fmt.Sprintf("http://localhost:%s/invoices/%s", port, pdfName)

	pdfFilename := filepath.Join(
		"invoices", pdfName,
	)

	err := os.WriteFile(pdfFilename, pdfData, 0644)
	if err != nil {
		return "", err
	}

	return pdfUrl, nil
}
