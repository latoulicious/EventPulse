# Use the official Golang image to create a build artifact.
FROM golang:1.23.5 AS builder

# Set the working directory inside the container.
WORKDIR /app

# Copy the Go module files.
COPY go.mod go.sum ./

# Download all the dependencies.
RUN go mod download

# Copy the rest of the source code.
COPY . .

# Build the application.
RUN CGO_ENABLED=0 GOOS=linux go build -o event-logger-go cmd/eventlogger/main.go

# Use a smaller base image for the final artifact.
FROM alpine:latest

# Set the working directory inside the container.
WORKDIR /app

# Copy the built application from the builder stage.
COPY --from=builder /app/event-logger-go .

# Expose the necessary ports.
EXPOSE 9092

# Run the application.
CMD ["./event-logger-go"]
