# Используем базовый образ с Go
FROM golang:1.23-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем все файлы проекта в контейнер
COPY . .

# Загружаем зависимости и собираем проект
RUN go mod tidy
RUN go build -o main main.go

# Используем легковесный образ для финального контейнера
FROM alpine:latest

# Устанавливаем зависимости
RUN apk --no-cache add ca-certificates

# Создаем рабочую директорию
WORKDIR /app

# Копируем собранное приложение из builder-образа
COPY --from=builder /app/main /app/main

# Копируем папку с шаблонами
COPY --from=builder /app/web /app/web

# Открываем порт (замените 8080 на нужный вам порт)
EXPOSE 8080

# Команда для запуска приложения
CMD ["/app/main"]
