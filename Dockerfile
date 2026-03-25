FROM golang:alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o api cmd/api/main.go

FROM alpine

WORKDIR /app

COPY --from=builder /app/api .

EXPOSE 8080

CMD ["./api"]