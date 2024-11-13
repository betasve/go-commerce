FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .


RUN go build -o gc-auth main.go

FROM alpine:latest

ENV APP_PORT=8001

COPY --from=builder /app/gc-auth /app/gc-auth
EXPOSE ${APP_PORT}

CMD ["/app/gc-auth"]
