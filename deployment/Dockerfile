# Используем официальный образ Go для сборки
FROM golang:1.23.4-alpine3.21

WORKDIR /app

# Копируем файлы проекта
COPY . .

# Загружаем зависимости
RUN go mod tidy

# Компилируем API и миграции
RUN go build -o api ./cmd/api/main.go
RUN go build -o migrate ./cmd/migration/main.go

# RUN apk --no-cache add ca-certificates

# Делаем файлы исполняемыми
RUN chmod +x /app/migrate /app/api

# Запускаем миграции, а затем основное приложение
CMD ["sh", "-c", "./migrate && ./api"]
