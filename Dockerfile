FROM golang:1.27-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY main.go ./

RUN go build -o server .

FROM alpine:3.22

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]