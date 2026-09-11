FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o inventory-service ./cmd/server

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/inventory-service .

EXPOSE 50051

CMD ["./inventory-service"]