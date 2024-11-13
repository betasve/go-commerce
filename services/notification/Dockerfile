FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .


RUN go build -o gc-notification main.go

FROM alpine:latest

ENV APP_PORT=8003

COPY --from=builder /app/gc-notification /app/gc-notification
EXPOSE ${APP_PORT}

CMD ["/app/gc-notification"]
