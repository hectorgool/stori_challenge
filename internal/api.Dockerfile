FROM golang:1.23-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the working directory for dependency installation
COPY go.mod go.sum ./

# Install Air for hot reloading during development
RUN go install github.com/air-verse/air@latest

# Copy the subdirectory containing the main.go file to the container
COPY cmd/app/ ./cmd/app/

# Install the project's dependencies
RUN go mod tidy

# Command to run the application using Air, specifying build command and output binary location
CMD ["air", "--build.cmd", "go build -o /app/tmp/main ./cmd/app", "--build.bin", "/app/tmp/main"]
