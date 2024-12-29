# Stage 1: Build the Go application
FROM golang:1.20 AS builder

# Set the working directory
WORKDIR /app

# Copy Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the application statically
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

FROM debian:bullseye-slim  

# Install CA certificates
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*


# Set working directory
WORKDIR /

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main .

# Expose the port
EXPOSE 8080

RUN apt-get update && apt-get install -y ca-certificates

# Run the binary
CMD ["/main"]


# Base image
# FROM golang:1.20

# # Set the working directory inside the container
# WORKDIR /app

# # Copy Go module files and download dependencies
# COPY go.mod go.sum ./
# RUN go mod download

# # Copy the rest of the application source code
# COPY . .

# # Build the application
# RUN go build -o main .

# # Expose the application port
# EXPOSE 8080

# # Run the application
# CMD ["./main"]
