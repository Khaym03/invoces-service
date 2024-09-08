# Usa una imagen base de Go
FROM golang:1.23.1-alpine AS go-builder

# Establece el directorio de trabajo en el contenedor
WORKDIR /app

# Copia los archivos go.mod y go.sum
COPY go.mod go.sum ./

# Descarga las dependencias de Go
RUN go mod download

# Copia el resto del código fuente de la aplicación
COPY . .

# Compila la aplicación Go
RUN go build -o bin cmd/api/main.go

# Segunda fase
FROM node:18-alpine AS node-builder

WORKDIR /app

COPY package*.json ./

RUN npm install

COPY . .

# Mezcla
FROM alpine:latest

# Instalar Chromium
RUN apk add --no-cache chromium

WORKDIR /app

COPY --from=go-builder /app/bin /app/bin
COPY --from=node-builder /app/node_modules /app/node_modules
COPY . .

ENV PORT="3000"

EXPOSE 3000

CMD ["./bin/main"]