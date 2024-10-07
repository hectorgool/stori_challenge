FROM golang:1.23-alpine

# Establecer directorio de trabajo
WORKDIR /app

# Instalar Air para hot reload desde el nuevo módulo
RUN go install github.com/air-verse/air@latest

# Copiar el proyecto al contenedor
COPY . .

# Instalar dependencias
RUN go mod tidy

# Ejecutar la aplicación con Air para hot reload
CMD ["air"]