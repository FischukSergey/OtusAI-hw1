# Stage 1: сборка Go-бинарника
FROM golang:1.26-alpine AS builder

WORKDIR /build

COPY backend/go.mod ./
RUN go mod download

COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# Stage 2: минимальный runtime-образ
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /build/server ./server
COPY frontend/ ./frontend/

ENV FRONTEND_DIR=/app/frontend
ENV PORT=8080

EXPOSE 8080

CMD ["./server"]
