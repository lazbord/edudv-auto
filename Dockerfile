# Step 1: Build the Go application in a builder stage
FROM golang:latest AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first, to cache the dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Set GOOS and GOARCH to ensure compatibility
ENV GOOS=linux
ENV GOARCH=amd64

# Build the Go app with CGO disabled to create a static binary
RUN CGO_ENABLED=0 go build -o main .

# Step 2: Create a smaller final image to run the Go application
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the compiled Go binary from the builder stage
COPY --from=builder /app/main .

# Copy the .env file into the container
COPY .env . 

# Ensure the binary has execution permissions
RUN chmod +x ./main

# Command to run the application
CMD ["./main"]
