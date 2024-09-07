run:
	@go run cmd/api/main.go

refresh:
	@templ generate

build-css:
	npx tailwindcss -i ./assets/css/styles.css -o ./assets/css/output.css

.PHONY: build-css


# Nuevo objetivo que incluye refresh, build-css y run
dev: refresh build-css run


build:
	go build -o bin cmd/api/main.go

