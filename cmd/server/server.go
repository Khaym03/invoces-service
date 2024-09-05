package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Khaym03/invoces-service/internal/components"
	"github.com/Khaym03/invoces-service/internal/models"
	"github.com/a-h/templ"
)

func main() {

	description := models.InvoiceDescription{
		Id:       123123,
		Date:     time.Now(),
		DueDate:  time.Now(),
		TotalDue: 50.00,
	}

	http.Handle("/", templ.Handler(components.Root(description)))
	http.HandleFunc("/x", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "generate/x.html")
	})

	fmt.Println("Listening on :3000")
	http.ListenAndServe(":3000", nil)
}
