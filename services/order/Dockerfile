FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .


RUN go build -o gc-order main.go

FROM alpine:latest

ENV APP_PORT=8004

COPY --from=builder /app/gc-order /app/gc-order
EXPOSE ${APP_PORT}

CMD ["/app/gc-order"]
