FROM golang:1.23-alpine

# Establecer el directorio de trabajo
WORKDIR /app

# Copiar go.mod y go.sum para instalar dependencias
COPY go.mod go.sum ./

# Instalar Air para hot reload
RUN go install github.com/air-verse/air@latest

# Copiar el subdirectorio que contiene el archivo main.go
COPY cmd/app/ ./cmd/app/

# Instalar dependencias
RUN go mod tidy

# Ejecutar la aplicación usando Air sin .air.toml
CMD ["air", "--build.cmd", "go build -o /app/tmp/main ./cmd/app", "--build.bin", "/app/tmp/main"]
