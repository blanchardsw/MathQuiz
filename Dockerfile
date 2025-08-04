# -------- Stage 1: Build --------
    FROM golang:1.21 AS builder

    WORKDIR /app
    
    # Copy go.mod and go.sum first to leverage caching
    COPY ./backend/go.mod ./backend/go.sum ./
    RUN go mod download
    
    # Copy the rest of the source code
    COPY ./backend ./
    
    # Build the binary
    RUN go build -mod=readonly -o main .
    
    # -------- Stage 2: Run --------
    FROM gcr.io/distroless/static:nonroot
    
    # Copy the binary from the builder stage
    COPY --from=builder /app/main /main
    
    # Run the binary
    CMD ["/main"]
    