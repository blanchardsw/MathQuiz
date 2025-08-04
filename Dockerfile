# -------- Stage 1: Build --------
    FROM golang:1.21 AS builder

    WORKDIR /app
    
    COPY ./backend/go.mod ./backend/go.sum ./
    RUN go mod download
    
    COPY ./backend ./
    RUN go build -mod=readonly -o main .
    
    # -------- Stage 2: Run --------
    FROM alpine:latest
    
    # Install shell and certs
    RUN apk --no-cache add ca-certificates bash
    
    # Copy binary and make it executable
    COPY --from=builder /app/main /main
    RUN chmod +x /main
    
    # Set working directory (optional)
    WORKDIR /
    
    # Explicitly set entrypoint
    ENTRYPOINT ["/main"]
    