FROM golang:1.21

WORKDIR /app

COPY ./backend/go.mod ./backend/go.sum ./
COPY ./backend ./backend

RUN go build -mod=readonly -o main ./backend

CMD ["./main"]