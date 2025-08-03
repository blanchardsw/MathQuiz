# Use official Go image
FROM golang:1.21

# Set working directory
WORKDIR /app

# Copy everything
COPY . .

# Build the backend
RUN go build -o main ./backend

# Run the binary
CMD ["./main"]