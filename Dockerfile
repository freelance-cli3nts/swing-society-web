# Start from golang base image
FROM golang:1.23-alpine AS builder

# Set working directory
WORKDIR /app

# Copy the entire project
COPY . .

# Set working directory to server
WORKDIR /app/server

# Download all dependencies
RUN go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Start a new stage from scratch
FROM alpine:latest  

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/server/main .

# Copy go.mod file for root directory detection
COPY --from=builder /app/go.mod .

# Copy static and template files from the root of the project
COPY --from=builder /app/static/ ./static/
COPY --from=builder /app/templates/ ./templates/

# # Copy config.json (if needed)
COPY --from=builder /app/config.json ./

# Set environment variable to indicate we're in a container
ENV DOCKER_CONTAINER=true

# Expose port 8080
EXPOSE 8080

# Command to run the executable
CMD ["./main"]