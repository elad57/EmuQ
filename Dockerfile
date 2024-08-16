# Stage 1: Build the Go binary
FROM golang:1.20-alpine AS builder

# Set environment variables
ENV GO111MODULE=on

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o main cmd/main.go

# Stage 2: Build a small image with the Go binary
FROM alpine:3.18

# Set environment variables
ENV GIN_MODE=release

# Create a non-root user
RUN adduser -D nonroot

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Change the ownership of the directory
RUN chown -R nonroot:nonroot /app

# Switch to the non-root user
USER nonroot

# Expose the application on port 8080
EXPOSE 8080 8082

# Command to run the executable
CMD ["./main"]
