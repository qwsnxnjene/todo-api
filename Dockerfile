# 1 этап — сборка
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Копируем go.mod и go.sum первыми (кэширование зависимостей)
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь исходный код
COPY . .

# Собираем статический бинарник (важно для alpine!)
RUN CGO_ENABLED=0 GOOS=linux go build -o todo-api .

# 2 этап — минимальный финальный образ
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем только бинарник из предыдущего этапа
COPY --from=builder /app/todo-api .

# Открываем порт
EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/tasks || exit 1

# Запускаем
CMD ["./todo-api"]