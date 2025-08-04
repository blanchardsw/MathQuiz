FROM golang:1.21

# Set working directory to backend (where go.mod lives)
WORKDIR /app/backend

# Copy go.mod and go.sum first
COPY ./backend/go.mod ./backend/go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the backend source code
COPY ./backend ./

# Build the Go binary
RUN go build -mod=readonly -o /app/main .

# Run the binary
CMD ["/app/main"]