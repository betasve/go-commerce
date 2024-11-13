FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .


RUN go build -o gc-product main.go

FROM alpine:latest

ENV APP_PORT=8005

COPY --from=builder /app/gc-product /app/gc-product
EXPOSE ${APP_PORT}

CMD ["/app/gc-product"]
