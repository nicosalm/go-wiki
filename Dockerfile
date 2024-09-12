
FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o wiki

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/wiki /app/wiki

COPY --from=builder /app/templates /app/templates

EXPOSE 8080

CMD ["./wiki"]

