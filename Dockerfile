# Stage 1: Build
FROM golang:1.26.4-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o task-api main.go

# Stage 2: Runtime
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/task-api .
EXPOSE 8080
CMD ["./task-api"]