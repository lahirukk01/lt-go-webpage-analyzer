# Use the official Golang image as the base image
FROM golang:1.23 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./

# Download and cache the dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
# RUN go build -o bin/app ./cmd
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/app ./cmd

# Use a minimal base image for the final image
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/bin/app .

# Make app executable
RUN chmod +x app

# Copy the public directory for static files
COPY --from=builder /app/public ./public

# Expose the port on which the application will run
EXPOSE 3000

# Command to run the application
CMD ["./app"]
