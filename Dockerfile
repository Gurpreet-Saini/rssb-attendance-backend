FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o attendance-api ./main.go

FROM alpine:3.21
WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/attendance-api ./attendance-api
COPY --from=builder /app/migrations ./migrations
EXPOSE 8080
CMD ["./attendance-api"]
