FROM golang:1.26 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /out/server ./server/cmd/


FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /out/server /usr/local/bin/server

CMD ["/usr/local/bin/server"]