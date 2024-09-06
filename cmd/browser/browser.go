package main

import (
	"fmt"
	"net/http"

	"github.com/Khaym03/invoces-service/internal/handlers"
)

func main() {
	handlers.ConfigureFileServers()

	http.HandleFunc("/generate", handlers.GenerateInvoiceHandler)

	fmt.Println("Listening on :3000")
	http.ListenAndServe(":3000", nil)

}

// w.Header().Set("Content-Type", "application/pdf")
// w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.pdf", filepath.Base(pdfURL)
