# Используем базовый образ с Go
FROM golang:1.23-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /OracleGo

COPY go.mod go.sum ./

RUN go mod download

# Копируем все файлы проекта в контейнер
COPY . .

# Загружаем зависимости и собираем проект
RUN go build -o main main.go

# Используем легковесный образ для финального контейнера
FROM python:3.10-slim

# Установка системных зависимостей
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential python3-dev libffi-dev gcc g++ && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

# Обновление pip и установка библиотек поэтапно
RUN pip install --no-cache-dir --upgrade pip setuptools wheel
RUN pip install --no-cache-dir pandas==1.5.3 numpy
RUN pip install --no-cache-dir torch torchvision torchaudio --extra-index-url https://download.pytorch.org/whl/cpu
RUN pip install --no-cache-dir lightautoml==0.3.8.1 nltk transformers

# Создаем рабочую директорию
WORKDIR /OracleGo

# Копируем собранное приложение из builder-образа
COPY --from=builder /OracleGo /OracleGo

# Команда для запуска приложения
CMD ["/OracleGo/main"]
