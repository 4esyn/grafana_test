# Используем официальный образ Go как базовый
FROM golang:1.22-alpine AS builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем файлы go.mod и go.sum и устанавливаем зависимости
COPY go.mod go.sum ./

# Устанавливаем swag
RUN go mod download && \
    go install github.com/swaggo/swag/cmd/swag@v1.8.1

# Копируем все файлы в рабочую директорию
COPY . .

# Генерируем документацию и собираем приложение
RUN swag init && \
    CGO_ENABLED=0 GOOS=linux go build -o main

# Используем официальный образ Alpine Linux как базовый
FROM alpine:3.14

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем бинарный файл из билдера в рабочую директорию
COPY --from=builder /app/main .
COPY --from=builder /app/docs ./docs

COPY start.sh .
RUN chmod +x start.sh

# Запускаем приложение
CMD ["./start.sh"]
