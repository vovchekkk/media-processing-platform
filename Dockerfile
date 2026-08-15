FROM golang:1.26 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/server ./cmd/server


FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/server ./server

CMD ["./server"]