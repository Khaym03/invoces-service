server:
	@go run cmd/server/server.go

pdf:
	@go run cmd/main.go

brow:
	@go run cmd/browser/browser.go

refresh:
	@templ generate

# Default target
all: refresh  pdf

.PHONY: all  refresh  pdf