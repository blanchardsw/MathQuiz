# -------- Stage 1: Build --------
    FROM golang:1.21 AS builder

    WORKDIR /app
    
    COPY ./backend/go.mod ./backend/go.sum ./
    RUN go mod download
    
    COPY ./backend ./
    RUN go build -mod=readonly -o main .
    
    # -------- Stage 2: Run --------
    FROM alpine:latest
    
    RUN apk --no-cache add ca-certificates
    
    COPY --from=builder /app/main /main
    RUN chmod +x /main
    
    ENTRYPOINT ["/main"]    