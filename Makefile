server:
	@go run cmd/server/server.go

pdf:
	@go run cmd/main.go

brow:
	@go run cmd/browser/browser.go

refresh:
	@templ generate

build-css:
	npx tailwindcss -i ./assets/css/styles.css -o ./assets/css/output.css

.PHONY: build-css

# Default target
all: refresh  pdf

.PHONY: all  refresh  pdf


# Nuevo objetivo que incluye refresh, build-css y brow
dev: refresh build-css brow


