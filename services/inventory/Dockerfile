FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .


RUN go build -o gc-inventory main.go

FROM alpine:latest

ENV APP_PORT=8002

COPY --from=builder /app/gc-inventory /app/gc-inventory
EXPOSE ${APP_PORT}

CMD ["/app/gc-inventory"]
