FROM golang:1.24.0-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o user_service cmd/main/main.go

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/user_service .
EXPOSE 8080
CMD ["./user_service"]