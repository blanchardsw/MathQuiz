# -------- Stage 1: Build --------
    FROM golang:1.21 AS builder

    WORKDIR /app
    
    # Copy Go module files and download dependencies
    COPY ./backend/go.mod ./backend/go.sum ./
    RUN go mod download
    
    # Copy source code and build statically for Alpine
    COPY ./backend ./
    RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=readonly -o /app/main .
    RUN ls -lh /app/main  # ✅ Debug: confirm binary exists
    
    # -------- Stage 2: Run --------
    FROM alpine:latest
    
    # Install certs and shell (optional)
    RUN apk --no-cache add ca-certificates bash
    
    # Copy binary from builder stage
    COPY --from=builder /app/main /main
    RUN chmod +x /main
    
    # Set working directory and entrypoint
    WORKDIR /
    ENTRYPOINT ["/main"]    