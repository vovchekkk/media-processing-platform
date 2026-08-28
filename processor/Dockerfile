FROM golang:1.26 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /out/processor ./processor/cmd/


FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /out/processor /usr/local/bin/processor

CMD ["/usr/local/bin/processor"]