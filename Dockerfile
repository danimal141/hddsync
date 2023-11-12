# Use the latest Go build environment
FROM golang:1.21-alpine as builder

# Set the working directory
WORKDIR /app

# Copy the Go module files
COPY go.mod .
COPY go.sum .

# Download Go module dependencies
RUN go mod download

# Copy the Go program and .env file
COPY hddsync.go .
COPY .env .

# Build the Go program
RUN go build -o sync

# Use a lightweight Alpine image for running the application
FROM alpine:latest

# Install necessary packages (rsync)
RUN apk add --no-cache rsync

# Copy the built executable from the builder stage
COPY --from=builder /app/sync /app/sync

# Set the working directory for the running container
WORKDIR /app

# Set the command to execute the sync program
CMD ["./sync"]
