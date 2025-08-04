# -------- Stage 1: Build --------
    FROM golang:1.21 AS builder

    WORKDIR /app
    
    COPY ./backend/go.mod ./backend/go.sum ./
    RUN go mod download
    
    COPY ./backend ./
    RUN go build -mod=readonly -o main .
    
    # -------- Stage 2: Run --------
    FROM alpine:latest
    
    # Install minimal dependencies (if needed)
    RUN apk --no-cache add ca-certificates
    
    # Copy the binary from builder
    COPY --from=builder /app/main /main
    
    # Set working directory (optional)
    WORKDIR /
    
    # Run the binary
    CMD ["/main"]
    