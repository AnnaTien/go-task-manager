# Stage 1: Build the Go application
# Use the official Go image as the build environment.
FROM golang:1.23-alpine AS builder

# Set the working directory inside the container.
WORKDIR /app

# Copy the source code into the container.
COPY . .

# Download and tidy up Go module dependencies.
RUN go mod tidy

# Build the Go application for a static executable.
# CGO_ENABLED=0 ensures the binary is statically linked, making it more portable.
RUN CGO_ENABLED=0 go build -o /go-task-manager ./cmd/main.go

# ---

# Stage 2: Create a lightweight final image
# Use a minimal Alpine image for the final production image.
FROM alpine:latest

# Set the working directory.
WORKDIR /root/

# Copy the statically linked binary from the builder stage.
COPY --from=builder /go-task-manager .

# Copy the configuration files.
COPY configs/ ./configs/

# Expose the port the application will run on.
EXPOSE 8080

# The command to run the application when the container starts.
CMD ["./go-task-manager", "api"]