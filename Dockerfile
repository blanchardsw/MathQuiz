# Use official Go image
FROM golang:1.21

# Set working directory inside container
WORKDIR /app

# Copy go.mod and go.sum first for caching
COPY ./backend/go.mod ./backend/go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the backend source code
COPY ./backend ./backend

# Build the Go binary
RUN go build -mod=readonly -o main ./backend

# Optional: expose port or set entrypoint
# EXPOSE 8080
# CMD ["./main"]