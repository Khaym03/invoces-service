package handlers

import (
	"net/http"
	"os"
	"path/filepath"
)

func ConfigureFileServers() {
	dir, _ := os.Getwd()

	http.Handle(
		"/assets/",
		http.StripPrefix(
			"/assets/",
			http.FileServer(
				http.Dir(filepath.Join(dir, "assets")),
			),
		),
	)

	http.Handle(
		"/generated/",
		http.StripPrefix(
			"/generated/",
			http.FileServer(
				http.Dir(filepath.Join(dir, "generated")),
			),
		),
	)

	http.Handle(
		"/generated-pdf/",
		http.StripPrefix(
			"/generated-pdf/",
			http.FileServer(
				http.Dir(filepath.Join(dir, "generated-pdf")),
			),
		),
	)
}
