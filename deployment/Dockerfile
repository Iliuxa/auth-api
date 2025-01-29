# Используем официальный образ Go для сборки
FROM golang:1.23.4 AS builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем файлы проекта в контейнер
COPY . .

# Загружаем зависимости
# RUN go mod tidy

# Компилируем приложение
RUN go build -o api ./cmd/api/main.go
RUN go build -o migrate ./cmd/migration/main.go

# Используем минимальный образ для продакшена
FROM alpine:latest

# Устанавливаем рабочую директорию
WORKDIR /app

# Устанавливаем зависимости (если нужно)
# RUN apk --no-cache add ca-certificates

# Копируем скомпилированные бинарники
COPY --from=builder /app/api .
COPY --from=builder /app/migrate .

# Запускаем миграции, а затем основное приложение
CMD ["sh", "-c", "./migrate && ./api"]
