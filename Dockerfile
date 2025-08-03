# Use official Go image
FROM golang:1.21

# Set working directory
WORKDIR /app

# Copy everything
COPY . .

# Build using readonly module mode
RUN go build -mod=readonly -o main ./backend

# Run the binary
CMD ["./main"]
