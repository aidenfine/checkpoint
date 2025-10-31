FROM golang:1.25-alpine AS builder
WORKDIR /app

RUN apk add --no-cache git

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/main .

RUN chmod +x ./main

EXPOSE 8080

CMD ["./main"]