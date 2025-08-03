FROM golang:1.21

WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker layer caching
COPY go.mod go.sum ./

# Copy the backend source code
COPY ./backend ./backend

# Build the binary
RUN go build -mod=readonly -o main ./backend

CMD ["./main"]
