FROM golang:1.22.6-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -v -o ./bin/dbo-test ./cmd/main.go

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/bin/dbo-test ./bin/dbo-test

EXPOSE 8080

CMD ["./bin/dbo-test"]
