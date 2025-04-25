# Use the official Golang image as the base image
FROM golang:1.24-alpine

# Set environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Install necessary dependencies
RUN apk add --no-cache git

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o main ./cmd/server

# Default values for host and port
ENV HOST=0.0.0.0
ENV PORT=8080

# Expose the port
EXPOSE ${PORT}

# Command to run the application with flags
CMD ["./main", "--host=${HOST}", "--port=${PORT}"]